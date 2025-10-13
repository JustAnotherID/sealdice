package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"

	"github.com/sealdice-ce/sealdice-ce/core/dice"
)

var (
	myDice *dice.Dice
	dm     *dice.DiceManager
)

func Bind(e *echo.Echo, _myDice *dice.DiceManager) {
	dm = _myDice
	myDice = _myDice.Dice[0]

	// 挂载 humaecho 到 echo 实例
	_ = humaecho.New(e, huma.DefaultConfig("Sealdiciapi", "1.0.0"))

	bindApiV1(e)
}

func bindApiV1(e *echo.Echo) {
	root := e.Group("/sd-api")

	bindInfoApi(root)
	bindLogApi(root)
	bindImConnectionApi(root)
	bindCustomTextApi(root)
	bindCustomReplyApi(root)
	bindSettingApi(root)
	bindPublicDiceApi(root)
	bindSignInApi(root)
	bindBackupApi(root)
	bindGroupApi(root)
	bindBanApi(root)
	bindDeckApi(root)
	bindJsApi(root)
	bindHelpDocApi(root)
	bindStoryApi(root)
	bindUtilsApi(root)
	bindCensorApi(root)
	bindResourceApi(root)
	bindVerifyApi(root)
	bindStoreApi(root)
}

func bindInfoApi(root *echo.Group) {
	info := root.Group("/info")
	info.GET("/preInfo", preInfo)
	info.GET("/baseInfo", baseInfo)
	info.GET("/hello", hello2)
}

func bindLogApi(root *echo.Group) {
	log := root.Group("/log")
	log.GET("/fetchAndClear", logFetchAndClear)
}

func bindImConnectionApi(root *echo.Group) {
	im := root.Group("/im_connections")
	im.GET("/list", ImConnections)
	im.GET("/get", ImConnectionsGet)
	im.POST("/qrcode", ImConnectionsQrcodeGet)
	im.POST("/gocq_captcha_set", ImConnectionsCaptchaSet)

	im.POST("/addGocq", ImConnectionsAddBuiltinGocq)
	im.POST("/addOnebot11ReverseWs", ImConnectionsAddReverseWs)
	im.POST("/addGocqSeparate", ImConnectionsAddGocqSeparate)
	im.POST("/addLagrange", ImConnectionsAddBuiltinLagrange)

	im.POST("/addDiscord", ImConnectionsAddDiscord)
	im.POST("/addKook", ImConnectionsAddKook)
	im.POST("/addTelegram", ImConnectionsAddTelegram)
	im.POST("/addMilky", ImConnectionsAddMilky)

	im.POST("/del", ImConnectionsDel)
	im.POST("/set_enable", ImConnectionsSetEnable)
	im.GET("/get_lgr_signinfo", ImConnectionsGetSignInfo)
	im.POST("/gocqhttpRelogin", ImConnectionsGocqhttpRelogin)
	im.GET("/gocq_config_download.zip", ImConnectionsGocqConfigDownload)
}

func bindCustomTextApi(root *echo.Group) {
	ct := root.Group("/configs/customText")
	ct.GET("/info", customTextInfo)
	ct.POST("/save", customTextSave)
	ct.POST("/preview-refresh", customTextPreviewRefresh)
}

func bindCustomReplyApi(root *echo.Group) {
	customReply := root.Group("/configs/custom_reply")
	customReply.GET("/info", customReplyGet)
	customReply.POST("/save", customReplySave)
	customReply.GET("/file_list", customReplyFileList)
	customReply.POST("/file_new", customReplyFileNew)
	customReply.POST("/file_delete", customReplyFileDelete)
	customReply.GET("/file_download", customReplyFileDownload)
	customReply.POST("/file_upload", customReplyFileUpload)
	customReply.GET("/debug_mode", customReplyDebugModeGet)
	customReply.POST("/debug_mode", customReplyDebugModeSet)
}

func bindSettingApi(root *echo.Group) {
	setting := root.Group("/dice")
	setting.GET("/config/get", DiceConfig)
	setting.POST("/config/set", DiceConfigSet)
	setting.GET("/config/advanced/get", DiceAdvancedConfigGet)
	setting.POST("/config/advanced/set", DiceAdvancedConfigSet)
	setting.POST("/config/mail_test", DiceMailTest)
	setting.POST("/exec", DiceExec)
	setting.GET("/recentMessage", DiceRecentMessage)
	setting.GET("/cmdList", DiceAllCommand)
	setting.POST("/upload_to_upgrade", DiceNewVersionUpload)
	setting.POST("/config/vm-version-set", vmVersionSet)
	setting.POST("/upgrade", upgrade)
	setting.POST("/force_stop", forceStop)
}

func bindPublicDiceApi(root *echo.Group) {
	publicDice := root.Group("/dice/public")
	publicDice.GET("/info", dicePublicInfo)
	publicDice.POST("/set", dicePublicSet)
}

func bindSignInApi(root *echo.Group) {
	signin := root.Group("/signin")
	signin.POST("/do", doSignIn)
	signin.GET("/salt", doSignInGetSalt)
	signin.GET("/checkSecurity", checkSecurity)
}

func bindBackupApi(root *echo.Group) {
	backup := root.Group("/backup")
	backup.GET("/list", backupGetList)
	backup.POST("/do_backup", backupExec)
	backup.GET("/config_get", backupConfigGet)
	backup.POST("/config_set", backupConfigSave)
	backup.GET("/download", backupDownload)
	backup.POST("/delete", backupDelete)
	backup.POST("/batch_delete", backupBatchDelete)
}

func bindGroupApi(root *echo.Group) {
	group := root.Group("/group")
	group.GET("/list", groupList)
	group.POST("/set_one", groupSetOne)
	group.POST("/quit_one", groupQuit)
}

func bindBanApi(root *echo.Group) {
	ban := root.Group("/banconfig")
	ban.GET("/list", banMapList)
	ban.GET("/get", banConfigGet)
	ban.POST("/set", banConfigSet)
	// ban.GET("/map_get", banMapGet)
	ban.POST("/map_delete_one", banMapDeleteOne)
	ban.POST("/map_add_one", banMapAddOne)
	// ban.POST("/map_set", banMapSet)
	ban.GET("/export", banExport)
	ban.POST("/import", banImport)
}

func bindDeckApi(root *echo.Group) {
	deck := root.Group("/deck")
	deck.GET("/list", deckList)
	deck.POST("/reload", deckReload)
	deck.POST("/upload", deckUpload)
	deck.POST("/enable", deckEnable)
	deck.POST("/delete", deckDelete)
	deck.POST("/check_update", deckCheckUpdate)
	deck.POST("/update", deckUpdate)
}

func bindJsApi(root *echo.Group) {
	js := root.Group("/js")
	js.POST("/reload", jsReload)
	js.POST("/execute", jsExec)
	js.POST("/upload", jsUpload)
	js.GET("/list", jsList)
	js.POST("/delete", jsDelete)
	js.GET("/get_record", jsGetRecord)
	js.POST("/shutdown", jsShutdown)
	js.GET("/status", jsStatus)
	js.POST("/enable", jsEnable)
	js.POST("/disable", jsDisable)
	js.POST("/check_update", jsCheckUpdate)
	js.POST("/update", jsUpdate)
	js.GET("/get_configs", handleGetConfigs)
	js.POST("/set_configs", handleSetConfigs)
	js.POST("/delete_unused_configs", handleDeleteUnusedConfigs)
	js.POST("/reset_config", handleResetConfig)
}

func bindHelpDocApi(root *echo.Group) {
	helpdoc := root.Group("/helpdoc")
	helpdoc.GET("/status", helpDocStatus)
	helpdoc.GET("/tree", helpDocTree)
	helpdoc.POST("/reload", helpDocReload)
	helpdoc.POST("/upload", helpDocUpload)
	helpdoc.POST("/delete", helpDocDelete)
	helpdoc.POST("/textitem/get_page", helpGetTextItemPage)
	helpdoc.GET("/config", helpGetConfig)
	helpdoc.POST("/config", helpSetConfig)
}

func bindStoryApi(root *echo.Group) {
	story := root.Group("/story")
	story.GET("/info", storyGetInfo)
	story.GET("/logs", storyGetLogs)
	story.GET("/logs/page", storyGetLogPage)
	story.GET("/items", storyGetItems)
	story.GET("/items/page", storyGetItemPage)
	story.DELETE("/log", storyDelLog)
	story.POST("/uploadLog", storyUploadLog)
	story.GET("/backup/list", storyGetLogBackupList)
	story.GET("/backup/download", storyDownloadLogBackup)
	story.POST("/backup/batch_delete", storyBatchDeleteLogBackup)
}

func bindUtilsApi(root *echo.Group) {
	utils := root.Group("/utils")
	utils.POST("/onebot", onebotTool)
	utils.GET("/ga/:uid", getGithubAvatar)
	utils.GET("/news", getNews)
	utils.POST("/check_news", checkNews)
	utils.GET("/get_token", getToken)
	utils.POST("/check_cron_expr", checkCronExpr)
	utils.GET("/check_network_health", checkNetworkHealth)
}

func bindCensorApi(root *echo.Group) {
	censor := root.Group("censor")
	censor.POST("/restart", censorRestart)
	censor.POST("/stop", censorStop)
	censor.GET("/status", censorGetStatus)
	censor.GET("/config", censorGetConfig)
	censor.POST("/config", censorSetConfig)
	censor.GET("/words", censorGetWords)
	censor.GET("/files", censorGetWordFiles)
	censor.POST("/files/upload", censorUploadWordFiles)
	censor.DELETE("/files", censorDeleteWordFiles)
	censor.GET("/files/template/toml", censorGetTomlFileTemplate)
	censor.GET("/files/template/txt", censorGetTxtFileTemplate)
	censor.GET("/logs/page", censorGetLogPage)
}

func bindResourceApi(root *echo.Group) {
	resource := root.Group("/resource")
	resource.GET("/page", resourceGetList)
	resource.GET("/download", resourceDownload)
	resource.POST("/upload", resourceUpload)
	resource.DELETE("/delete", resourceDelete)
	resource.GET("/data", resourceGetData)
}

func bindVerifyApi(root *echo.Group) {
	verify := root.Group("/verify")
	verify.GET("/generate_code", verifyGenerateCode)
}

func bindStoreApi(root *echo.Group) {
	store := root.Group("/store")
	store.GET("/backend/list", storeBackendList)
	store.POST("/backend/add", storeAddBackend)
	store.DELETE("/backend/remove", storeRemoveBackend)
	store.GET("/recommend", storeRecommend)
	store.GET("/page", storeGetPage)
	store.POST("/download", storeDownload)
	store.POST("/rating", storeRating)
}
