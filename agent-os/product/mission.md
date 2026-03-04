# Product Mission

## Problem

Automatisera kontinuerlig förbättring och underhåll av mjukvaruprojekt genom agenter som arbetar dygnet runt och skapar pull requests för att implementera förbättringar, fixa buggar och skriva tester.

## Target Users

Projektägare och utvecklingsteam som vill skala kodunderhåll, snabba upp feedback-loopar och avlasta rutinuppgifter.

## Solution

Issue-driven multi-agent-system som automatiskt:

- Skapar en agent per GitHub-issue och orkestrerar flera agenter samtidigt.
- Automatiserar worktree/branch-creation, implementation, test- och refactor-loopar.
- Hanterar modell-provider-failover och kontextbevarande mellan providers.
- Skapar ephemeral miljöer för manuell testning innan PR merge.
- Följer strikta säkerhetsprinciper: agenter körs med separata konton och isolerade arbetsmiljöer, utan åtkomst till användarens personliga maskin.
