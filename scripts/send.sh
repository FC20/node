#!/usr/bin/env bash
# Sends coins from one test account to another
ouroboroscli tx send  --gas auto --gas-adjustment 1.5 $(ouroboroscli keys show jack -a) $(ouroboroscli keys show alice -a) 2000000000ouro