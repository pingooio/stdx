package autoload

/*
	You can just read the .env file on import just by doing

		import _ "github.com/pingooio/stdx/dotenv/autoload"

	And bob's your mother's brother
*/

import dotenv "github.com/pingooio/stdx/dotenv"

func init() {
	dotenv.Load()
}
