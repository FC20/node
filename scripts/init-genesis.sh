#!/usr/bin/env bash
# Sets up a network from scratch

PASSWORD="abcdef123"

rm -rf ~/.ouroborosd/
rm -rf ~/.ouroboroscli/

# Initializes blockchain
ouroborosd init anonymous --chain-id ouroboros

# Creates a few testing accounts
echo ${PASSWORD} | ouroboroscli keys add jack
echo ${PASSWORD} | ouroboroscli keys add alice

# Adds genesis account
ouroborosd add-genesis-account $(ouroboroscli keys show jack -a) 10000000000000ouro,10000000000000stake

# Конфигурируем cli
ouroboroscli config chain-id ouroboros
ouroboroscli config output json
ouroboroscli config indent true
ouroboroscli config trust-node true

# Creates genesis tx
echo ${PASSWORD} | ouroborosd gentx --name jack --amount 100000000stake

ouroborosd collect-gentxs