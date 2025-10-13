package api

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/sealdice-ce/sealdice-ce/core/dice"
)

func ImConnections(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	return c.JSON(http.StatusOK, myDice.ImSession.EndPoints)
}

func ImConnectionsGet(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	v := struct {
		ID string `form:"id" json:"id"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		for _, i := range myDice.ImSession.EndPoints {
			if i.ID == v.ID {
				return c.JSON(http.StatusOK, i)
			}
		}
	}
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsSetEnable(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		ID     string `form:"id"     json:"id"`
		Enable bool   `form:"enable" json:"enable"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		for _, i := range myDice.ImSession.EndPoints {
			if i.ID == v.ID {
				i.SetEnable(myDice, v.Enable)
				return c.JSON(http.StatusOK, i)
			}
		}
	}

	myDice.LastUpdatedTime = time.Now().Unix()
	myDice.Save(false)
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsGetSignInfo(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	data, err := dice.LagrangeGetSignInfo(myDice)
	if err != nil {
		return Error(&c, "读取SignInfo失败", Response{})
	}
	return Success(&c, Response{
		"data": data,
	})
}

func ImConnectionsDel(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		ID string `form:"id" json:"id"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		for index, i := range myDice.ImSession.EndPoints {
			if i.ID == v.ID {
				// 禁用该endpoint防止出问题
				i.SetEnable(myDice, false)
				// 待删除的EPInfo落库，保留其统计数据
				i.StatsDump(myDice)
				// TODO: 注意 这个好像很不科学
				// i.diceServing = false
				switch i.Platform {
				case "QQ":
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					if i.ProtocolType == "onebot" {
						pa := i.Adapter.(*dice.PlatformAdapterGocq)
						if pa.BuiltinMode == "lagrange" || pa.BuiltinMode == "lagrange-gocq" {
							dice.BuiltinQQServeProcessKillBase(myDice, i, true)
							// 经测试，若不延时，可能导致清理对应目录失败（原因：文件被占用）
							time.Sleep(1 * time.Second)
							dice.LagrangeServeRemoveConfig(myDice, i)
						} else {
							dice.BuiltinQQServeProcessKill(myDice, i)
						}
					}
					return c.JSON(http.StatusOK, i)
				case "DISCORD":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "KOOK":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "TG":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "MC":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "DODO":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "DINGTALK":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "SLACK":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "SEALCHAT":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				}
			}
		}
		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Save(false)
	}
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsQrcodeGet(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	v := struct {
		ID string `form:"id" json:"id"`
	}{}
	err := c.Bind(&v)
	if err != nil {
		return c.JSON(http.StatusNotFound, nil)
	}

	for _, i := range myDice.ImSession.EndPoints {
		if i.ID != v.ID {
			continue
		}
		switch i.ProtocolType {
		case "onebot", "":
			pa := i.Adapter.(*dice.PlatformAdapterGocq)
			if pa.GoCqhttpState == dice.StateCodeInLoginQrCode {
				return c.JSON(http.StatusOK, map[string]string{
					"img": "data:image/png;base64," + base64.StdEncoding.EncodeToString(pa.GoCqhttpQrcodeData),
				})
			}
		}
		return c.JSON(http.StatusOK, i)
	}
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsCaptchaSet(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	v := struct {
		ID   string `form:"id"   json:"id"`
		Code string `form:"code" json:"code"`
	}{}
	err := c.Bind(&v)
	if err != nil {
		return err
	}

	for _, i := range myDice.ImSession.EndPoints {
		if i.ID == v.ID {
			switch i.ProtocolType {
			case "onebot", "":
				pa := i.Adapter.(*dice.PlatformAdapterGocq)
				if pa.GoCqhttpState == dice.GoCqhttpStateCodeInLoginBar {
					pa.GoCqhttpLoginCaptcha = v.Code
					return c.String(http.StatusOK, "")
				}
			}
		}
	}
	return c.String(http.StatusNotFound, "")
}

func ImConnectionsGocqhttpRelogin(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		ID string `form:"id" json:"id"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		for _, i := range myDice.ImSession.EndPoints {
			if i.ID == v.ID {
				myDice.Logger.Warnf("relogin %s", v.ID)
				i.Adapter.DoRelogin()
				return c.JSON(http.StatusOK, nil)
			}
		}
	}
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsGocqConfigDownload(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	id := c.QueryParam("id")
	for _, i := range myDice.ImSession.EndPoints {
		if i.ID == id {
			buf := packGocqConfig(i.RelWorkDir)
			return c.Blob(http.StatusOK, "", buf.Bytes())
		}
	}

	return c.String(http.StatusNotFound, "")
}

type AddDiscordEcho struct {
	Token              string
	ProxyURL           string
	ReverseProxyUrl    string
	ReverseProxyCDNUrl string
}

func ImConnectionsAddDiscord(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := &AddDiscordEcho{}
	err := c.Bind(&v)
	if err == nil {
		conn := dice.NewDiscordConnItem(dice.AddDiscordEcho(*v))
		conn.Session = myDice.ImSession
		pa := conn.Adapter.(*dice.PlatformAdapterDiscord)
		pa.Session = myDice.ImSession
		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Save(false)
		go dice.ServeDiscord(myDice, conn)
		return c.JSON(http.StatusOK, conn)
	}
	return c.String(430, "") // 这个是非标的，呃。。
}

func ImConnectionsAddKook(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		Token string `json:"token" yaml:"token"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		conn := dice.NewKookConnItem(v.Token)
		conn.Session = myDice.ImSession
		pa := conn.Adapter.(*dice.PlatformAdapterKook)
		pa.Session = myDice.ImSession
		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Save(false)
		go dice.ServeKook(myDice, conn)
		return c.JSON(http.StatusOK, conn)
	}
	return c.String(430, "")
}

func ImConnectionsAddTelegram(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		Token    string `json:"token"    yaml:"token"`
		ProxyURL string `json:"proxyURL" yaml:"proxyURL"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		conn := dice.NewTelegramConnItem(v.Token, v.ProxyURL)
		conn.Session = myDice.ImSession

		// myDice.Logger.Infof("成功创建endpoint")
		pa := conn.Adapter.(*dice.PlatformAdapterTelegram)
		pa.Session = myDice.ImSession
		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Save(false)
		go dice.ServeTelegram(myDice, conn)
		return c.JSON(http.StatusOK, conn)
	}
	return c.String(430, "")
}

func ImConnectionsAddMilky(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		WsGateway   string `json:"wsGateway"   yaml:"wsGateway"`
		RestGateway string `json:"restGateway" yaml:"restGateway"`
		Token       string `json:"token"       yaml:"token"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		conn := dice.NewMilkyConnItem(dice.AddMilkyEcho{
			Token:       v.Token,
			WsGateway:   v.WsGateway,
			RestGateway: v.RestGateway,
		})
		pa := conn.Adapter.(*dice.PlatformAdapterMilky)
		pa.Session = myDice.ImSession
		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Save(false)
		go dice.ServeMilky(myDice, conn)
		return c.JSON(http.StatusOK, conn)
	}
	return c.String(430, "")
}

func ImConnectionsAddBuiltinGocq(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		Account       string `json:"account"       yaml:"account"`
		Password      string `json:"password"      yaml:"password"`
		Protocol      int    `json:"protocol"`
		AppVersion    string `json:"appVersion"`
		UseSignServer bool   `json:"useSignServer"`
	}{}

	err := c.Bind(&v)
	if err == nil {
		uid := v.Account
		if checkUidExists(c, uid) {
			return nil
		}

		conn := dice.NewGoCqhttpConnectInfoItem(v.Account)
		conn.UserID = dice.FormatDiceIDQQ(uid)
		conn.Session = myDice.ImSession
		pa := conn.Adapter.(*dice.PlatformAdapterGocq)
		pa.InPackGoCqhttpProtocol = v.Protocol
		pa.InPackGoCqhttpPassword = v.Password
		pa.InPackGoCqhttpAppVersion = v.AppVersion
		pa.Session = myDice.ImSession
		pa.UseSignServer = v.UseSignServer

		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
		myDice.LastUpdatedTime = time.Now().Unix()

		dice.GoCqhttpServe(myDice, conn, dice.GoCqhttpLoginInfo{
			Password:      v.Password,
			Protocol:      v.Protocol,
			AppVersion:    v.AppVersion,
			IsAsyncRun:    true,
			UseSignServer: v.UseSignServer,
		})
		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Save(false)
		return c.JSON(http.StatusOK, conn)
	}
	return c.String(430, "")
}

func ImConnectionsAddGocqSeparate(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		Account     string `json:"account"     yaml:"account"`
		ConnectURL  string `json:"connectUrl"  yaml:"connectUrl"`  // 连接地址
		AccessToken string `json:"accessToken" yaml:"accessToken"` // 访问令牌
	}{}

	err := c.Bind(&v)
	if err == nil {
		uid := v.Account
		if checkUidExists(c, uid) {
			return nil
		}

		conn := dice.NewGoCqhttpConnectInfoItem("")
		conn.UserID = dice.FormatDiceIDQQ(uid)
		conn.Session = myDice.ImSession

		pa := conn.Adapter.(*dice.PlatformAdapterGocq)
		pa.Session = myDice.ImSession

		// 三项设置
		conn.RelWorkDir = "x" // 此选项已无意义
		pa.ConnectURL = v.ConnectURL
		pa.AccessToken = v.AccessToken

		pa.UseInPackClient = false

		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
		conn.SetEnable(myDice, true)

		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Save(false)
		return c.JSON(http.StatusOK, conn)
	}
	return c.String(430, "")
}

func ImConnectionsAddReverseWs(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		Account     string `json:"account"     yaml:"account"`
		ReverseAddr string `json:"reverseAddr" yaml:"reverseAddr"`
	}{}

	err := c.Bind(&v)
	if err == nil {
		uid := v.Account
		if checkUidExists(c, uid) {
			return nil
		}

		conn := dice.NewGoCqhttpConnectInfoItem(v.Account)
		conn.UserID = dice.FormatDiceIDQQ(uid)
		conn.Session = myDice.ImSession

		pa := conn.Adapter.(*dice.PlatformAdapterGocq)
		pa.Session = myDice.ImSession

		pa.IsReverse = true
		pa.ReverseAddr = v.ReverseAddr

		pa.UseInPackClient = false

		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
		conn.SetEnable(myDice, true)

		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Save(false)
		return c.JSON(http.StatusOK, conn)
	}
	return c.String(430, "")
}

func ImConnectionsAddBuiltinLagrange(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, Response{"testMode": true})
	}

	v := struct {
		Account           string `json:"account"           yaml:"account"`
		SignServerName    string `json:"signServerName"    yaml:"signServerName"`
		SignServerVersion string `json:"signServerVersion" yaml:"signServerVersion"`
		IsGocq            bool   `json:"isGocq"            yaml:"isGocq"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		uid := v.Account
		if checkUidExists(c, uid) {
			return nil
		}

		conn := dice.NewLagrangeConnectInfoItem(v.Account, v.IsGocq)
		conn.UserID = dice.FormatDiceIDQQ(uid)
		conn.Session = myDice.ImSession
		pa := conn.Adapter.(*dice.PlatformAdapterGocq)
		// pa.InPackGoCqhttpProtocol = v.Protocol
		pa.Session = myDice.ImSession

		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
		myDice.LastUpdatedTime = time.Now().Unix()
		uin, err := strconv.ParseInt(v.Account, 10, 64)
		if err != nil {
			return err
		}
		pa.SignServerName = v.SignServerName
		pa.SignServerVer = v.SignServerVersion
		dice.LagrangeServe(myDice, conn, dice.LagrangeLoginInfo{
			UIN:               uin,
			SignServerName:    v.SignServerName,
			SignServerVersion: v.SignServerVersion,
			IsAsyncRun:        true,
		})
		return c.JSON(http.StatusOK, v)
	}

	return c.String(430, "")
}
