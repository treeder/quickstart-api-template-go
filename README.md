# FLAP API Template

just testing rapid app creation

[frontend template](https://github.com/treeder/flap-ui-template)

* Update go mod, change module name at top to new repo github.com/treeder/xxx
* Create firebase project
* Upgrade to pay as you go
* Create firestore
* Go to settings -> Service accounts, click "Generate new private key". This will download a JSON file. 
* Run `base64 -w 0 account.json` to get encoded version of the file (for secrets)
* Make `.env` file with G_KEY and G_PROJECT_ID
    * And or add this into your repos codespace secrets as G_KEY and add G_PROJECT_ID
* make run (boom)
