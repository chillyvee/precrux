# Precrux

Configuring horcux can be challenging.  Even for adminstrators who are good at their jobs, reconfiguring new horcrux clusters can be time consuming.

Precrux allows a controlling node (SNITCH) to configure remote signing servers (CHASER).

For security, copy and paste the certificate from CHASER to the SNITCH to estalibsh a secure channel.

SHUTDOWN THE CHASER on the remote signers once horcrux is configured.  Horcrux runs without the CHASER.


# Prepartion for Local and Remote Horcrux nodes

Install precrux and horcrux on each system (SNITCH and 3x CHASER)

Keep your existing validator running

Make a directory to contain precrux configuration.  For example

```
CHAIN=uni
mkdir $HOME/.precrux
```


# Prepare each chaser

Assume you have 3 chasers named "red", "green" and "blue"


# On remote (chaser) nodes

Start the chaser on the remote node.  Specify the PORT to listen for incoming precrux information.

NOTE: Incoming port should be open on the firewall, but exposure should be limited because precrux will be shutdown after configuration receipt.

NOTE: Close firewall port for extra security

```
precrux chaser start red --port 5050
```


# On control (snitch) node

Register the chaser locally
```
precrux remote add red
```

# Prepare a chain for horcrux.  

For example "uni" which is Juno's testnet

```
CHAIN=uni
mkdir $HOME/.precrux/$CHAIN
```

Copy the priv_validator_key.json from your existing validator to the local computer

```
mkdir $HOME/.precrux/$CHAIN
cp priv_validator_key.json $HOME/.precrux/$CHAIN
```

Backup your priv_validator_key.json from your existing validator 

Rename or delete the priv_validator_key.json from your existing validator 

Copy precrux.yaml into the chain configuration directory


```
CHAIN=uni
cp precrux.yaml $HOME/.precrux/$CHAIN
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
  # name - A single word naming the remote signer.  Example: red, west-1, horcrux-a
  # p2p-listen - tcp://ip:port for remote-signers to talk to each other
  # debug-addr - ip:port for prometheus metrics (accessible at ip:port/metrics)
  # priv-val-addr - list of sentries that apply only to this signer 
  #                 tcp://ip:port of sentry/validator (port should match config.toml priv_validator_laddr = "0.0.0.0:4000")
  #                 Also possible to list multiple values comma separated: tcp://chain-node-1:1234,tcp://chain-node-2:1234
  -
    name: red
    p2p-listen: tcp://1.1.1.1:2001
    priv-val-addr: tcp://5.5.5.5:4001
    debug-addr: 0.0.0.0:3001
    chaser-addr: 1.1.1.1:5050
  -
    name: green
    p2p-listen: tcp://2.2.2.2:2001
    debug-addr: 0.0.0.0:3001
    priv-val-addr: tcp://6.6.6.6:4001
  -
    name: blue
    p2p-listen: tcp://3.3.3.3:2001
    debug-addr: 0.0.0.0:3001
    priv-val-addr: tcp://7.7.7.7:4001
```

# Locally Generate horcrux files for the chain

```
precrux generate uni
```

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
