package coindaemons

func MineBlock() error {
	for _, cd := range CoinDaemons {
		err := cd.MineBlocks(1)
		if err != nil {
			return err
		}
	}
	return nil
}
