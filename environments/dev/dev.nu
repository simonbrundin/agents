#!/usr/bin/env nu
cd /home/simon/repos/software-agent/environments/dev

# -----------------------------------------------
# Aktivera Mise-paket
# -----------------------------------------------

def activate-mise [] {
  let mise_path = (^mise env | parse -r 'export PATH=\'(.+)\'' | get capture0.0 | split row (char esep))
  $env.PATH = ($mise_path | append $env.PATH)
}

activate-mise



def get-local-ip [] {
    try {
        let ip_output = (ip route get 1.1.1.1 | complete)
        if $ip_output.exit_code == 0 {
            let ip = ($ip_output.stdout | str trim | parse -r 'src\s+(\S+)' | get capture0.0)
            $ip
        } else {
            let ips = (hostname -I | str trim | split row ' ')
            $ips.0
        }
    } catch {
        "127.0.0.1"
    }
}

# -----------------------------------------------
# Schema Sync - Håll dev-databas i synk med schema.sql
# -----------------------------------------------

def get-schema-hash [file: string] {
    ^sha256sum $file | str trim | split row " " | get 0
}

let schema_file = "/home/simon/repos/software-agent/environments/prod/db/schema.sql"
let schema_hash_file = "/home/simon/repos/software-agent/environments/dev/.schema.hash"

let current_hash = (get-schema-hash $schema_file)
let saved_hash = try { (open $schema_hash_file | str trim) } catch { "" }

if $current_hash != $saved_hash {
    print "Schema ändrat - synkar med schema.sql..."

    try {
        docker exec software-agent-dev-postgres psql -U postgres -d software_agent -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
        docker exec software-agent-dev-postgres psql -U postgres -d software_agent -f /docker-entrypoint-initdb.d/01-schema.sql
        docker exec software-agent-dev-postgres psql -U postgres -d software_agent -f /docker-entrypoint-initdb.d/02-mock-data.sql
    } catch {
        print "Kör docker compose up -d postgres först för att initiera databasen"
    }

    try {
        hasura metadata apply --endpoint http://localhost:8080 --admin-secret hasura-dev-secret
    } catch {
        print "Kör hasura först (kör dev igen efter att tilt startat)"
    }

    $current_hash | save -f $schema_hash_file
    print "Schema synkat!"
}

# -----------------------------------------------
# Variables
# -----------------------------------------------

let local_ip = (get-local-ip)
let namespace = "software-agent-dev"
let link = $"http://($local_ip):3000"

# -----------------------------------------------
# Skicka notis till telefon
# -----------------------------------------------

curl -d $link ntfy.sh/simonbrundin-dev-notification

# -----------------------------------------------
# Kör Tilt
# -----------------------------------------------

try {
  let tilt_pids = (pgrep tilt | lines | where $it != "")
  if ($tilt_pids | length) > 0 {
    print $"Avslutar befintliga Tilt-processer: ($tilt_pids)"
    pkill -9 tilt
  }
} catch {
}

try {
  fuser -k 10350/tcp
} catch {
}

# -----------------------------------------------
# Välj Tilt mode med fzf
# -----------------------------------------------

print "\n🚀 Välj utvecklingsmiljö:\n"

let modes = [
  "local - Docker Compose (standard)",
  "kubernetes - Full cluster miljö"
]

let selection = try {
  ($modes | fzf --prompt="Mode: " --height=40% --reverse | str trim)
} catch {
  "local - Docker Compose (standard)"
}

let mode = if ($selection | str contains "kubernetes") {
  "kubernetes"
} else {
  "local"
}

print $"\n✓ Startar Tilt i ($mode) mode...\n"

if $mode == "kubernetes" {
  tilt up --namespace $namespace -- mode=kubernetes
} else {
  tilt up --namespace $namespace -- mode=local
}
