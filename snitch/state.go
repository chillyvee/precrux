package snitch

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/strangelove-ventures/horcrux/signer"
	tmBytes "github.com/tendermint/tendermint/libs/bytes"
	tmjson "github.com/tendermint/tendermint/libs/json"
)

// Snippet Taken from https://raw.githubusercontent.com/tendermint/tendermint/main/privval/file.go
// FilePVLastSignState stores the mutable part of PrivValidator.
type FilePVLastSignState struct {
	Height int64 `json:"height"`
	Round  int32 `json:"round"`
	Step   int8  `json:"step"`
}

// From horcrux/cmd/horcrux/cmd/sign_state.go
// SignState stores signing information for high level watermark management.
type SignState struct {
	Height          int64            `json:"height"`
	Round           int64            `json:"round"`
	Step            int8             `json:"step"`
	EphemeralPublic []byte           `json:"ephemeral_public"`
	Signature       []byte           `json:"signature,omitempty"`
	SignBytes       tmBytes.HexBytes `json:"signbytes,omitempty"`
	cache           map[signer.HRSKey]signer.SignStateConsensus

	filePath string
}

func (s *Snitch) ImportPrivValidatorState() (*FilePVLastSignState, error) {
	// Allow user to paste in priv_validator_state.json
	fmt.Println("")
	fmt.Println("********************************************************************************")
	fmt.Println("IMPORTANT: Your validator should already be STOPPED.  You must copy the latest state..")
	fmt.Println("********************************************************************************")
	fmt.Println("")
	fmt.Println("Paste your old priv_validator_state.json.  Input a blank line after the pasted JSON to continue.")
	fmt.Println("")
	fmt.Println("\tNOTE: Type YOLO will skip importing your signing state.  This greatly increases your risk of double signing and a permanent tombstone.")
	fmt.Println("")

	var textBuffer strings.Builder

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			break
		}
		textBuffer.WriteString(scanner.Text())
	}
	finalJSON := textBuffer.String()

	pvState := &FilePVLastSignState{}
	var err error

	// Unsafe Branch
	if strings.TrimSpace(finalJSON) == "YOLO" {
		fmt.Println("********  DANGER  **************************************************************")
		fmt.Println("You have chosen to allow double signing")
		fmt.Println("This may cause permanent unrecoverable tombstoning for your validator")
		fmt.Println("********************************************************************************")
		fmt.Println("")
		fmt.Printf("Run 'precrux generate %s' again to properly import state", s.ChainName)
		fmt.Println("")
		fmt.Println("Think about your life choices for 5 seconds")
		fmt.Println("")
		<-time.After(5 * time.Second)
		pvState.Height = 0
		pvState.Round = 0
		pvState.Step = 0
		return pvState, nil
	}

	// Parse JSON and return pvState
	err = tmjson.Unmarshal([]byte(finalJSON), &pvState)
	if err != nil {
		fmt.Println("Error parsing priv_validator_state.json")
		return nil, err
	}

	fmt.Printf("Imported Priv Val Sign State: \n"+
		"  Height:    %v\n"+
		"  Round:     %v\n"+
		"  Step:      %v\n",
		pvState.Height, pvState.Round, pvState.Step)
	return pvState, nil
}

func (s *Snitch) ImportAndWriteSignState() error {
	var err error

	// Import State
	pvState := &FilePVLastSignState{}
	for {
		pvState, err = s.ImportPrivValidatorState()
		if err == nil {
			break
		}
	}

	signState := signer.SignStateConsensus{
		Height:    pvState.Height,
		Round:     int64(pvState.Round),
		Step:      pvState.Step,
		Signature: nil,
		SignBytes: nil,
	}

	// Create Cosigner State
	for _, co := range s.Config.Cosigners {
		// Create folder
		dir := s.ChaserStateDir(co)
		fmt.Printf("Creating %s\n", dir)
		err := os.MkdirAll(dir, 0700)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}

		// Recreate privValStateFile if necessary
		pv, err := signer.LoadOrCreateSignState(s.ChaserFile(co, s.ChaserStatePvsFile()))
		if err != nil {
			log.Fatal(err)
			return err
		}

		// shareStateFile does not exist during default config init, so create if necessary
		share, err := signer.LoadOrCreateSignState(s.ChaserFile(co, s.ChaserStateSssFile()))
		if err != nil {
			log.Fatal(err)
			return err
		}
		pv.EphemeralPublic, share.EphemeralPublic = nil, nil

		fmt.Printf("Saving New Sign State: \n"+
			"  Height:    %v\n"+
			"  Round:     %v\n"+
			"  Step:      %v\n",
			signState.Height, signState.Round, signState.Step)

		// func (signState *SignState) Save(ssc SignStateConsensus, lock *sync.Mutex, async bool) error {
		err = pv.Save(signState, nil, false)
		if err != nil {
			fmt.Printf("error saving privval sign state : %+v", err)
			return err
		}
		err = share.Save(signState, nil, false)
		if err != nil {
			fmt.Printf("error saving share sign state : %+v", err)
			return err
		}

		fmt.Printf("Update Successful\n")
	}
	return nil
}
