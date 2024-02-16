// Copyright (C) 2024, AllianceBlock. All rights reserved.
// See the file LICENSE for licensing terms.

package actions

var (
	OutputValueZero    = []byte("value is zero")
	OutputMemoTooLarge = []byte("memo is too large")

	OutputLockupPeriodInvalid       = []byte("lockup period is invalid")
	OutputStakeMissing              = []byte("stake is missing")
	OutputUnauthorized              = []byte("unauthorized")
	OutputInvalidNodeID             = []byte("invalid node ID")
	OutputDifferentNodeIDThanStaked = []byte("node ID is different than staked")
	OutputAssetMissing              = []byte("asset missing")
	OutputSymbolEmpty               = []byte("symbol is empty")
	OutputSymbolTooLarge            = []byte("symbol is too large")
	OutputDecimalsTooLarge          = []byte("decimal is too large")
	OutputMetadataEmpty             = []byte("metadata is empty")
	OutputMetadataTooLarge          = []byte("metadata is too large")
	OutputNotWarpAsset              = []byte("not warp asset")
	OutputWrongDestination          = []byte("wrong destination")
	OutputWarpAsset                 = []byte("warp asset")
	OutputAnycast                   = []byte("anycast output")

	// import_asset.go
	OutputConflictingAsset       = []byte("warp has same asset as another")
	OutputSymbolIncorrect        = []byte("symbol is incorrect")
	OutputDecimalsIncorrect      = []byte("decimal is incorrect")
	OutputWarpVerificationFailed = []byte("warp verification failed")
	OutputInvalidDestination     = []byte("invalid destination")
	OutputMustFill               = []byte("must fill request")

	// mint_asset.go
	OutputAssetIsNative = []byte("cannot mint native asset")
	OutputWrongOwner    = []byte("wrong owner")

	// register_validator_stake.go
	OutputDifferentSignerThanActor   = []byte("different signer than actor")
	OutputStakedAmountInvalid        = []byte("staked amount must be between 1.5 million and 1 billion")
	OutputInvalidStakeStartTime      = []byte("invalid stake start time")
	OutputInvalidStakeEndTime        = []byte("invalid stake end time")
	OutputInvalidStakeDuration       = []byte("invalid stake duration")
	OutputInvalidDelegationFeeRate   = []byte("delegation fee rate must be over 2 and under 100")
	OutputValidatorAlreadyRegistered = []byte("validator already registered")
)
