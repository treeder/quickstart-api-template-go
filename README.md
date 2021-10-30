# FLAP API Template

just testing rapid app creation

frontend: https://github.com/treeder/temp-ui

* Update go mod, change module name at top to new repo github.com/treeder/xxx
* Create firebase project
* Upgrade to pay as you go
* Create firestore
* Go to settings -> Service accounts, click "Generate new private key". This will download a JSON file. 
* Run `base64 -w 0 account.json` to get encoded version of the file (for secrets)
* Add this into repo codespace secrets as G_KEY and add G_PROJECT_ID
* Might need to restart codespace to pick up secrets
* make run (boom)
