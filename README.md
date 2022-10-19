# Status of project

Precrux is pre-alpha and is in active development.  This tool should be verified against testnets before attempting to apply to mainnets

# Precrux

Configuring horcux can be challenging.  Even for adminstrators who are good at their jobs, reconfiguring new horcrux clusters can be time consuming.

Precrux allows a controlling node (SNITCH) to configure remote signing servers (CHASER).

SHUTDOWN THE CHASER on the remote signers once horcrux is configured.  Horcrux runs without the CHASER.


# Security

Security and management of any key material is outside the scope of this service. Always consider your own security and risk profile when dealing with sensitive keys, services, or infrastructure.

# No Liability

This software comes as is, without any warranty or condition, and no contributor will be liable to anyone for any damages related to this software or this license, under any kind of legal claim.

# Who should / should NOT use this software

You should be fairly comfortable running tendermint/cosmos-sdk chains before attempting to use horcrux.  There are known cases of teams double signing for a permanent tombstone during this process.

We advise against using this configuration tool during a chain halt.  Fix your chain first before attempting any configuration.

Teams who often need to request help from others during setup or upgrade of the chain itself are likely the teams who should NOT be using this software.

Do not use this in production until you have learned how to use it properly in testnet.  If you need help accessing a testnet, contact us or any chain admin.

# Important before configuration:

* Your validator node should be shutdown
* priv_validator_key.json file should be backed up
* priv_validator_key.json should be removed from the validator/sentry

If you are more knowledgeable, you can keep running your validator node, but there are documented cases of teams double signing for a permanent tombstone during similar processes.

# Prepartion for Local and Remote Horcrux nodes

Install precrux and horcrux on each system (SNITCH and 3x CHASER)

Keep your existing validator running

Make a directory to contain precrux configuration.  For example

```
CHAIN=uni
mkdir $HOME/.precrux
```

Any other directory is acceptable as long as you run precrux from that directory.  For example:

```
cd $HOME/.precrux
```


# Prepare each chaser

Assume you have 3 chasers named "red", "green" and "blue"


# On remote (CHASER) nodes

Start the chaser on the remote node.  Specify the PORT to listen for incoming precrux information.

NOTE: Incoming port should be open on the firewall, but exposure should be limited because precrux will be shutdown after configuration receipt.

NOTE: Close firewall port for extra security

```
cd $HOME/.precrux
precrux chaser start red --port 5050
```

Repeat for chasers named "blue" and "green"

A certificate will print on the screen for you to copy and paste into the local SNITCH.

# On local control (SNITCH) node

Register the chaser locally
```
cd $HOME/.precrux
precrux remote add red IP:PORT
```


When prompted, paste the certificate printed upon chaser start.

Repeat for chsers named "blue" and "green"

# Prepare a chain for horcrux.  

For example "uni" which is Juno's testnet

```
cd $HOME/.precrux
CHAIN=uni
mkdir $HOME/.precrux/$CHAIN
```

Copy the priv_validator_key.json from your existing validator to the local computer

```
cp priv_validator_key.json $HOME/.precrux/$CHAIN/priv_validator_key.json 
```

Backup your priv_validator_key.json from your existing validator 

Rename or delete the priv_validator_key.json from your existing validator 

Copy precrux.yaml into the chain configuration directory


```
CHAIN=uni
cp precrux.yaml $HOME/.precrux/$CHAIN/precrux.yaml
```

Edit the configuration file to describe all horcrux nodes

```
# chain-name - Recommend using a single word.  Will be used for directory creation
chain-name: uni

# chain-id - Must match chain-id specified in genesis.json for the target chain
chain-id: uni-5

# threshold - Typically 2 (for 2-of-3 signing)
threshold: 2

# shares - Typically 3 (for 2-of-3 signing)
shares: 3

# rpc-timeout: Typically "1500ms"
rpc-timeout: 1500ms

# cosigner - Configure remote signers
cosigners:
  # By convention, share IDs are assigned 1,2,3 to the signers below in the order written
  # chaser-name - A single word naming the remote signer.  Example: red, west-1, horcrux-a
  # p2p-listen - tcp://ip:port for remote-signers to talk to each other
  # debug-addr - ip:port for prometheus metrics (accessible at ip:port/metrics)
  # priv-val-addr - list of sentries that apply only to this signer 
  #                 tcp://ip:port of sentry/validator (port should match config.toml priv_validator_laddr = "0.0.0.0:4000")
  #                 Also possible to list multiple values comma separated: tcp://chain-node-1:1234,tcp://chain-node-2:1234
  -
    chaser-name: red
    p2p-listen: tcp://1.1.1.1:2001
    debug-addr: 0.0.0.0:3001
    priv-val-addr: tcp://5.5.5.5:4001
  -
    chaser-name: green
    p2p-listen: tcp://2.2.2.2:2001
    debug-addr: 0.0.0.0:3001
    priv-val-addr: tcp://6.6.6.6:4001
  -
    chaser-name: blue
    p2p-listen: tcp://3.3.3.3:2001
    debug-addr: 0.0.0.0:3001
    priv-val-addr: tcp://7.7.7.7:4001
```

# Locally Generate horcrux files for the chain

Your prepared files should look like this
```
└── uni
    ├── precrux.yaml
    └── priv_validator_key.json
```

Generate all the configuration files for your remote signers

```
precrux generate uni
```

When prompted, copy and paste the contents of .chaindirectory/data/priv_validator_state.json

Enter a blank line to continue

Horcrux uses this information to prevent double signing.



# Push configuration for the chain to the remote chaser

```
precrux push uni red
precrux push uni green
precrux push uni blue
```

# Push additional chains from local to remote chaser

```
precrux generate gaia 
precrux push gaia red
precrux push gaia green
precrux push gaia blue
```

# Where files arrive

Files will be written to the current directory under the chain name.

For example:

```
$HOME/.precrux/uni
```

# Shut down precrux on the remote nodes

Type Control+C to terminate precrux

# Start horcrux on each node

```
horcrux cosigner start --home $HOME/.precrux/uni
```

# Remove key files

Saving precrux.yaml may be helpful for regenerating threshold signing configuration if you need to create new remote signers.

The cosigner diretory ($HOME/.precrux/uni/cosigner) can be deleted from the local computer if the remote horcrux is operating properly


