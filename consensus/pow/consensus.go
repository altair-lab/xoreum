package pow

import ()

// implement Engine's functions


// VerifyHeader checks whether a header conforms to the consensus rules of the
// stock Ethereum ethash engine.
func (pow *Pow) VerifyHeader() error {

}

// VerifySeal implements consensus.Engine, checking whether the given block satisfies
// the PoW difficulty requirements.
func (pow *Pow) VerifySeal() error {
}
