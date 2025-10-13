<template>
  <!-- <div style="position: relative;"> -->
  <Teleport to="#root">
    <div style="position: absolute; right: 40px; bottom: 60px">
      <n-button type="primary" class="btn-add" circle @click="addOne">
        <template #icon>
          <i-carbon-add-large />
        </template>
      </n-button>
    </div>
  </Teleport>
  <!-- </div> -->

  <div v-if="!store.curDice.conns || (store.curDice.conns && store.curDice.conns.length === 0)">
    <n-text>似乎还没有账号，</n-text>
    <n-button type="primary" size="small" text link @click="addOne">点我添加一个</n-button>
  </div>

  <n-list hoverable style="display: flex; flex-wrap: wrap">
    <n-list-item
      v-for="(i, index) in reactive(store.curDice.conns)"
      :key="index"
      style="min-width: 20rem; flex: 1 0 50%; grow: 0">
      <n-thing :bordered="false" content-class="py-4">
        <template #header>
          <span style="word-break: break-all">
            {{ i.nickname || '<"未知">' }}({{ i.userId }})
          </span>
        </template>
        <template #header-extra>
          <n-button size="small" type="error" secondary @click="doRemove(i)">
            <template #icon>
              <i-carbon-row-delete />
            </template>
            删除
          </n-button>
        </template>

        <div
          v-if="
            i.adapter?.loginState === goCqHttpStateCode.InLoginQrCode && store.curDice.qrcodes[i.id]
          "
          style="position: absolute; width: 17rem; height: 14rem; background: #fff">
          <div style="margin-left: 2rem">需要同账号的手机 QQ 扫码登录 (限 2 分钟内完成):</div>
          <img
            style="
              margin-left: -3rem;
              image-rendering: pixelated;
              width: 10rem;
              height: 10rem;
              margin-left: 3.5rem;
              margin-top: 2rem;
            "
            :src="store.curDice.qrcodes[i.id]" />
        </div>

        <n-descriptions ref="formRef" :model="i" :column="1" label-placement="left">
          <n-descriptions-item label="状态">
            <n-flex>
              <div v-if="i.state === 0">
                <n-tag size="small" type="error" :bordered="false">断开</n-tag>
              </div>
              <div v-if="i.state === 2">
                <n-tag size="small" type="warning" :bordered="false">连接中</n-tag>
              </div>
              <div v-if="i.state === 1">
                <n-tag size="small" type="success" :bordered="false">已连接</n-tag>
              </div>
              <div v-if="i.state === 3">
                <n-tag size="small" type="error" :bordered="false">失败</n-tag>
              </div>
              <n-tooltip
                v-if="
                  Math.round(new Date().getTime() / 1000) -
                    i.adapter?.inPackGoCqHttpLastRestricted <
                  30 * 60
                ">
                <template #trigger>
                  <n-tag size="small" type="warning">风控</n-tag>
                </template>
                {{
                  `看到这个标签是因为最近 20 分钟内有风控警告，将在重新登录后临时消除。触发时间：` +
                  dayjs.unix(i.adapter?.inPackGoCqHttpLastRestricted).fromNow()
                }}
              </n-tooltip>
            </n-flex>
          </n-descriptions-item>

          <n-descriptions-item label="群组数量">
            <div>{{ i.groupNum }}</div>
          </n-descriptions-item>

          <n-descriptions-item label="累计响应指令">
            <div>{{ i.cmdExecutedNum }}</div>
          </n-descriptions-item>

          <n-descriptions-item label="上次执行指令">
            <div v-if="i.cmdExecutedLastTime">
              {{ dayjs.unix(i.cmdExecutedLastTime).fromNow() }}
            </div>
            <div v-else>从未</div>
          </n-descriptions-item>

          <n-descriptions-item v-if="i.adapter?.connectUrl" label="连接地址">
            <!-- <n-input v-model:value="i.connectUrl"></n-input> -->
            <div>{{ i.adapter?.connectUrl }}</div>
          </n-descriptions-item>

          <n-descriptions-item v-if="i.adapter?.isReverse" label="服务地址">
            <div>{{ i.adapter?.reverseAddr }}/ws</div>
          </n-descriptions-item>

          <template
            v-if="
              i.platform === 'QQ' && i.protocolType === 'onebot' && i.adapter.builtinMode === 'gocq'
            ">
            <n-descriptions-item v-if="i.adapter.useInPackGoCqhttp" label="协议">
              <div v-if="i.adapter?.inPackGoCqHttpProtocol === 0">Unset</div>
              <div v-if="i.adapter?.inPackGoCqHttpProtocol === 1">Android</div>
              <div v-if="i.adapter?.inPackGoCqHttpProtocol === 2">Android 手表</div>
              <div v-if="i.adapter?.inPackGoCqHttpProtocol === 3">MacOS</div>
              <div v-if="i.adapter?.inPackGoCqHttpProtocol === 5">iPad</div>
              <div v-if="i.adapter?.inPackGoCqHttpProtocol === 6">AndroidPad</div>
            </n-descriptions-item>
            <n-descriptions-item v-if="i.adapter.useInPackGoCqhttp" label="协议版本">
              <div v-if="i.adapter?.inPackGoCqHttpAppVersion === ''">未指定</div>
              <div
                v-if="
                  i.adapter?.inPackGoCqHttpAppVersion && i.adapter.inPackGoCqHttpAppVersion !== ''
                ">
                {{ i.adapter.inPackGoCqHttpAppVersion }}
              </div>
            </n-descriptions-item>
            <n-descriptions-item v-if="i.adapter.useInPackGoCqhttp" label="协议实现">
              <div v-if="i.adapter?.implementation === 'gocq' || i.adapter?.implementation === ''">
                Go-Cqhttp
              </div>
            </n-descriptions-item>
            <n-descriptions-item v-else-if="i.adapter?.isReverse" label="特殊">
              <div>反向 WS</div>
            </n-descriptions-item>
            <n-descriptions-item v-else label="特殊">
              <div>分离部署</div>
            </n-descriptions-item>
          </template>

          <template
            v-if="
              i.platform === 'QQ' &&
              i.protocolType === 'onebot' &&
              i.adapter.builtinMode === 'lagrange'
            ">
            <n-descriptions-item label="接入方式">
              <div>内置客户端</div>
            </n-descriptions-item>
            <n-descriptions-item label="签名版本">
              <div>{{ i.adapter.signServerVer }}</div>
            </n-descriptions-item>
            <n-descriptions-item label="签名服务">
              <div>{{ i.adapter.signServerName }}</div>
            </n-descriptions-item>
          </template>

          <template
            v-if="
              i.platform === 'QQ' &&
              i.protocolType === 'onebot' &&
              i.adapter.builtinMode === 'lagrange-gocq'
            ">
            <n-descriptions-item label="接入方式">
              <div>内置 gocq</div>
            </n-descriptions-item>
            <n-descriptions-item label="签名版本">
              <div>{{ i.adapter.signServerVer }}</div>
            </n-descriptions-item>
            <n-descriptions-item label="签名服务">
              <div>{{ i.adapter.signServerName }}</div>
            </n-descriptions-item>
          </template>

          <template
            v-if="
              i.platform === 'QQ' && i.protocolType === 'onebot' && i.adapter.builtinMode === 'gocq'
            ">
            <n-descriptions-item label="其他">
              <n-tooltip placement="top-start">
                <template #trigger>
                  <n-button type="primary" @click="doGocqExport(i)">导出</n-button>
                </template>
                导出 gocq 设置，用于转分离部署
              </n-tooltip>
            </n-descriptions-item>
          </template>
        </n-descriptions>

        <template
          #action
          v-if="
            ![goCqHttpStateCode.InLogin, goCqHttpStateCode.InLoginQrCode].includes(
              i.adapter?.loginState,
            )
          ">
          <n-button-group>
            <n-tooltip placement="bottom-start">
              <template #trigger>
                <n-button type="warning" @click="askGocqhttpReLogin(i)">重新登录</n-button>
              </template>
              如果日志中出现帐号被风控，可以试试这个功能
            </n-tooltip>
            <n-tooltip placement="bottom-start">
              <template #trigger>
                <n-button v-if="i.enable" type="warning" @click="askSetEnable(i, false)">
                  禁用
                </n-button>
                <n-button v-else type="warning" @click="askSetEnable(i, true)">启用</n-button>
              </template>
              离线/启用此账号，重启骰子后仍将保持离线/启用状态
            </n-tooltip>
          </n-button-group>
        </template>
      </n-thing>
    </n-list-item>
  </n-list>
  <n-modal
    v-model:show="dialogFormVisible"
    title="帐号登录"
    :mask-closable="false"
    :close-on-esc="false"
    :closeable="false"
    preset="card">
    <template #header-extra>
      <n-button tertiary @click="openSocks"> 辅助工具 -13325 端口 </n-button>
    </template>
    <template v-if="form.step === 1">
      <n-alert
        v-if="
          store.diceServers.length > 0 &&
          store.diceServers[0].baseInfo.containerMode &&
          form.accountType === 15
        "
        type="warning"
        :closable="false"
        class="mb-6">
        当前为容器模式，内置客户端被禁用。
      </n-alert>
      <n-alert
        v-if="
          store.diceServers.length > 0 &&
          store.diceServers[0].baseInfo.containerMode &&
          form.accountType === 16
        "
        type="warning"
        :closable="false"
        class="mb-6">
        当前为容器模式，内置 gocq 被禁用。
      </n-alert>

      <n-form :model="form">
        <n-form-item label="账号类型" :label-width="formLabelWidth">
          <n-select v-model:value="form.accountType" :options="accountTypes" />
        </n-form-item>

        <n-form-item
          v-if="form.accountType === 15 || form.accountType === 16"
          label="账号"
          :label-width="formLabelWidth"
          required>
          <n-input v-model:value="form.account"></n-input>
        </n-form-item>
        <n-form-item
          v-if="form.accountType === 15 || form.accountType === 16"
          label="签名版本"
          :label-width="formLabelWidth"
          required>
          <template #label>
            <span>
              签名版本
              <n-tooltip>
                <template #trigger>
                  <n-icon><i-carbon-help-filled /></n-icon>
                </template>
                如果不知道这是什么，请保持默认选中
              </n-tooltip>
            </span>
          </template>

          <n-spin size="small" :show="!signInfoLoaded">
            <n-radio-group
              v-model:value="form.signServerVersion"
              :disabled="!signInfoLoaded"
              @change="signServerVersionChange">
              <template v-for="info in signInfos">
                <n-radio
                  v-if="!info.ignored"
                  :key="info.version"
                  :label="info.version"
                  :value="info.version">
                  <n-flex size="small">
                    <n-tag v-if="info.selected" size="small" type="success" :bordered="false">
                      最新
                    </n-tag>
                    {{ info.version }}
                  </n-flex>
                </n-radio>
              </template>
              <n-radio key="custom" label="自定义" value="自定义"></n-radio>
              <template #label>
                <template v-for="info in signInfos">
                  <div
                    :key="info.version"
                    v-if="info.version === form.signServerVersion && info.selected"
                    style="display: flex; align-items: center">
                    <span style="float: left; margin-right: 0.5rem">
                      {{ form.signServerVersion }}
                    </span>
                    <n-tag size="small" type="success" :bordered="false">最新</n-tag>
                  </div>
                </template>
              </template>
            </n-radio-group>
            <n-text v-if="signVerWarningText !== ''" type="warning">
              {{ signVerWarningText }}
            </n-text>
          </n-spin>
        </n-form-item>
        <n-form-item
          v-if="form.accountType === 15 || form.accountType === 16"
          label="签名服务"
          required>
          <template #label>
            <span>
              签名服务
              <n-tooltip>
                <template #trigger>
                  <n-icon><i-carbon-help-filled /></n-icon>
                </template>
                如果不知道这是什么，请保持默认选中
              </n-tooltip>
            </span>
          </template>

          <n-spin size="small" :show="!signInfoLoaded">
            <n-radio-group
              v-if="form.signServerVersion !== '自定义'"
              v-model:value="form.signServerName"
              :disabled="!signInfoLoaded"
              @change="signServerServerChange">
              <!-- placeholder -->
              <n-radio v-if="!signInfoLoaded">海豹 V2</n-radio>

              <template v-for="info in signInfos">
                <template v-if="info.version === form.signServerVersion && !info.ignored">
                  <template v-for="server in info.servers">
                    <n-radio
                      v-if="!server.ignored"
                      :key="server.name"
                      :label="server.name"
                      :value="server.name">
                      <n-flex align="center">
                        <span>{{ server.name }}</span>
                        <n-tag
                          v-if="server.latency < 120"
                          size="small"
                          type="success"
                          :bordered="false">
                          {{ server.latency }} ms
                        </n-tag>
                        <n-tag
                          v-else-if="server.latency >= 120 && server.latency < 360"
                          size="small"
                          :bordered="false"
                          type="warning">
                          {{ server.latency }} ms
                        </n-tag>
                        <n-tag
                          v-else-if="server.latency >= 360"
                          type="error"
                          size="small"
                          :bordered="false">
                          {{ server.latency }} ms
                        </n-tag>
                      </n-flex>
                    </n-radio>
                  </template>
                </template>
              </template>
            </n-radio-group>
            <n-input
              v-else
              v-model:value="form.signServerName"
              placeholder="请输入自定义签名地址"></n-input>
            <n-text v-if="signServerWarningText !== ''" type="warning" size="small">
              {{ signServerWarningText }}
            </n-text>
          </n-spin>
        </n-form-item>

        <n-form-item
          v-if="form.accountType === 6"
          label="账号"
          :label-width="formLabelWidth"
          required>
          <n-input v-model:value="form.account"></n-input>
        </n-form-item>

        <n-form-item
          v-if="form.accountType === 6"
          label="连接地址"
          :label-width="formLabelWidth"
          required>
          <n-input
            v-model:value="form.connectUrl"
            placeholder="正向WS连接地址，如 ws://127.0.0.1:1234"
            type="text"></n-input>
        </n-form-item>
        <n-form-item v-if="form.accountType === 6" label="访问令牌" :label-width="formLabelWidth">
          <n-input
            v-model:value="form.accessToken"
            placeholder="gocqhttp配置的access token，没有不用填写"></n-input>
        </n-form-item>

        <n-form-item
          v-if="form.accountType === 11"
          label="账号"
          :label-width="formLabelWidth"
          required>
          <n-input v-model:value="form.account"></n-input>
        </n-form-item>
        <n-form-item
          v-if="form.accountType === 11"
          label="连接地址"
          :label-width="formLabelWidth"
          required>
          <n-input
            v-model:value="form.reverseAddr"
            placeholder="反向WS服务地址，如 0.0.0.0:4001 (允许全部IP连入，4001端口)"
            type="text"></n-input>
        </n-form-item>

        <n-form-item
          v-if="form.accountType === 13"
          label="连接地址"
          :label-width="formLabelWidth"
          required>
          <n-input
            v-model:value="form.url"
            placeholder="连接地址，如 ws://127.0.0.1:3212/ws/seal"
            type="text"></n-input>
        </n-form-item>
        <n-form-item
          v-if="form.accountType === 13"
          label="Token"
          :label-width="formLabelWidth"
          required>
          <n-input
            v-model:value="form.token"
            type="text"
            placeholder="填入平台管理界面中获取的token"></n-input>
        </n-form-item>

        <n-form-item
          v-if="form.accountType === 14"
          label="平台"
          :label-width="formLabelWidth"
          required>
          <n-radio-group v-model:value="form.platform">
            <n-radio-button value="QQ" />
          </n-radio-group>
        </n-form-item>
        <n-form-item
          v-if="form.accountType === 14"
          label="主机"
          :label-width="formLabelWidth"
          required>
          <n-input
            v-model:value="form.host"
            placeholder="Satori 服务的地址，如 127.0.0.1"></n-input>
        </n-form-item>
        <n-form-item
          v-if="form.accountType === 14"
          label="端口"
          :label-width="formLabelWidth"
          required>
          <n-input-number v-model:value="form.port as any" placeholder="如 5500" />
        </n-form-item>
        <n-form-item v-if="form.accountType === 14" label="Token" :label-width="formLabelWidth">
          <n-input
            v-model:value="form.token"
            placeholder="填入鉴权 token，没有时无需填写"></n-input>
        </n-form-item>

        <n-form-item
          v-if="form.accountType === 1"
          label="Token"
          :label-width="formLabelWidth"
          required>
          <n-input v-model:value="form.token"></n-input>
          <small>
            <div>提示：首先去 discord 开发者平台创建一个新的 Application</div>
            <div>https://discord.com/developers/applications</div>
            <div>点击 New Application 创建之后进入应用，然后点 bot，Add bot</div>
            <div>然后把 Privileged Gateway Intents 下面的三个开关打开</div>
            <div>最后把 bot 的 token 复制下来粘贴进来</div>
          </small>
        </n-form-item>
        <n-form-item
          v-if="form.accountType === 1"
          label="http 代理地址"
          :label-width="formLabelWidth">
          <n-input v-model:value="form.proxyURL" placeholder="例：http://127.0.0.1:7890" />
        </n-form-item>
        <n-form-item
          v-if="form.accountType === 1"
          label="反向代理地址"
          :label-width="formLabelWidth">
          <n-input
            v-model:value="form.reverseProxyUrl"
            placeholder="此地址需要代理到 https://discord.com/ 通常来说正向代理和反向代理只需要一个" />
          <n-input
            v-model:value="form.reverseProxyCDNUrl"
            placeholder="此地址需要代理到 https://cdn.discordapp.com/ " />
          <div style="color: #aa4422">
            注意：反向代理是全局生效的，你一旦设置反向代理地址，这个骰子的所有 Discord
            连接都会使用这个地址
          </div>
        </n-form-item>

        <n-form-item
          v-if="form.accountType === 2"
          label="Token"
          :label-width="formLabelWidth"
          required>
          <n-input v-model:value="form.token" />
          <small>
            <div>提示：进入 KOOK 开发者平台创建一个新的应用</div>
            <div>https://developer.kookapp.cn/app/index</div>
            <div>点击新建应用 创建之后进入应用，然后点机器人</div>
            <div>把机器人的 token 复制下来粘贴进来</div>
          </small>
        </n-form-item>

        <n-form-item
          v-if="form.accountType === 3"
          label="Token"
          :label-width="formLabelWidth"
          required>
          <n-input v-model:value="form.token" />
          <small>
            <div>提示：私聊 BotFather(https://t.me/BotFather)</div>
            <div>使用/newbot 申请一个新的机器人</div>
            <div>
              按照指示创建机器人之后，在 Bot setting 里面把 Group privacy 里面 privacy mode 关掉
            </div>
            <div>把机器人的 token 复制下来粘贴进来</div>
          </small>
        </n-form-item>
        <n-form-item
          v-if="form.accountType === 3"
          label="http 代理地址"
          :label-width="formLabelWidth">
          <n-input v-model:value="form.proxyURL" placeholder="http://127.0.0.1:7890" />
        </n-form-item>
      </n-form>

      <n-form-item v-if="form.accountType === 17" label="Token" :label-width="formLabelWidth">
        <n-input v-model:value="form.token" />
      </n-form-item>
      <n-form-item
        v-if="form.accountType === 17"
        label="Websocket Gateway"
        :label-width="formLabelWidth"
        required>
        <n-input v-model:value="form.wsGateway" placeholder="ws://127.0.0.1:3000/event"></n-input>
      </n-form-item>
      <n-form-item
        v-if="form.accountType === 17"
        label="REST Gateway"
        :label-width="formLabelWidth"
        required>
        <n-input
          v-model:value="form.restGateway"
          placeholder="http://127.0.0.1:3000 (注意，不要加上/api的后缀)"></n-input>
      </n-form-item>
    </template>
    <template v-else-if="form.step === 2">
      <n-timeline style="min-height: 260px">
        <n-timeline-item
          v-for="(activity, index) in activities"
          :key="index"
          :type="activity.type as any"
          :color="activity.color"
          :size="activity.size as any"
          :hollow="activity.hollow">
          <div>{{ activity.content }}</div>
          <div v-if="index === 2 && isTestMode">
            <div>欢迎体验海豹骰点核心，展示模式下不提供登录功能，请谅解。</div>
            <div>如需测试指令，请移步“指令测试”界面。</div>
            <div>此外，数据会定期自动重置。</div>
            <div>展示版本未必是最新版，建议您下载体验。</div>
            <n-button style="margin-top: 1rem" @click="formClose">再会</n-button>
          </div>
          <div
            v-else-if="
              index === 2 && curConn.adapter?.loginState === goCqHttpStateCode.InLoginQrCode
            ">
            <div>登录需要扫码验证，请使用登录此账号的手机 QQ 扫描二维码以继续登录：</div>
            <img
              :src="store.curDice.qrcodes[curConn.id]"
              style="width: 20rem; height: 20rem; image-rendering: pixelated" />
          </div>

          <div
            v-else-if="
              index === 2 && curConn.adapter?.loginState === goCqHttpStateCode.LoginFailed
            ">
            <div>
              <div>登录失败！可能是以下原因：</div>
              <ul>
                <li>密码错误 (注意检查大小写)</li>
                <li>二维码获取失败 (日志中会有“二维码错误”)</li>
                <li>扫二维码超时过期 (约 2 分钟)</li>
                <li>海豹发生了故障</li>
              </ul>
              <div>具体请参见日志。如果不出现“确定”按钮，请直接刷新。</div>
              <div>若删除账号重复尝试无果，请回报 bug 给开发者。</div>
              <n-button text link href="javascript:window.location.reload()" type="primary">
                点我前往日志界面
              </n-button>
            </div>
          </div>
        </n-timeline-item>
      </n-timeline>
    </template>
    <template v-else-if="form.step === 3">
      <n-result
        icon="success"
        title="成功啦!"
        sub-title="现在账号状态应该是“已连接”了，去试一试骰子吧！">
        <!-- <template #extra></template> -->
      </n-result>
    </template>
    <template v-else-if="form.step === 4">
      <n-result icon="success" title="成功啦!" sub-title="操作完成，现在可以进行下一步了">
        <!-- <template #extra></template> -->
      </n-result>
    </template>

    <template #footer>
      <n-flex size="small" class="dialog-footer">
        <template v-if="form.step === 1">
          <n-button @click="dialogFormVisible = false">取消</n-button>
          <n-button
            type="primary"
            :disabled="
              ((form.accountType === 1 || form.accountType === 2 || form.accountType === 3) &&
                form.token === '') ||
              (form.accountType === 6 && (form.account === '' || form.connectUrl === '')) ||
              (form.accountType === 11 && (form.account === '' || form.reverseAddr === '')) ||
              ((form.accountType === 15 || form.accountType === 16) &&
                (form.account === '' ||
                  form.signServerVersion === '' ||
                  form.signServerName === '')) ||
              (form.accountType === 17 && (form.wsGateway === '' || form.restGateway === ''))
            "
            @click="goStepTwo">
            下一步
          </n-button>
        </template>
        <template v-if="form.isEnd">
          <n-button @click="formClose">确定</n-button>
        </template>
      </n-flex>
    </template>
  </n-modal>

  <div
    v-show="dialogSlideVisible"
    id="slide"
    style="
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(1, 1, 1, 0.7);
    ">
    <iframe
      id="slideIframe"
      ref="slideIframe"
      src="about:blank"
      rel="noreferrer"
      style="width: 100%; height: 100%"></iframe>
    <div
      v-show="slideBottomShow"
      style="
        position: absolute;
        bottom: 0;
        width: 100%;
        height: 100px;
        display: flex;
        justify-content: center;
        flex-direction: column;
        align-items: center;
      ">
      <div style="margin-bottom: 0.5rem">
        <a style="line-break: anywhere; font-size: 0.5rem" :href="slideLink" target="_blank"
          >方式 2:新页面打开 (如无法验证)</a
        >
      </div>
      <n-button secondary type="primary" @click="closeCaptchaFrame">关闭，滑条完成后点击</n-button>
    </div>
  </div>
</template>

<script lang="tsx" setup>
import { reactive } from 'vue';
import { useDialog, useMessage } from 'naive-ui';
import { useStore, goCqHttpStateCode } from '~/store';
import type { DiceConnection } from '~/store';
import { sleep } from '~/utils';
import * as dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { urlBase } from '~/backend';
import {
  getLagrangeSignInfo,
  postConnectionDel,
  postConnectionQrcode,
  postConnectSetEnable,
  type SignInfo,
} from '~/api/v1/im_connections';
import { postToolOnebot } from '~/api/v1/others';

dayjs.extend(relativeTime);

const fullActivities = [
  {
    content: '正在生成虚拟设备信息',
    size: 'large',
    type: 'primary',
    hollow: true,
  },
  {
    content: '正在生成登录配置文件',
    size: 'large',
    color: '#0bbd87',
    hollow: true,
  },
  {
    content: '进行登录……',
    size: 'large',
    flag: true,
  },
  {
    content: '完成！',
    type: 'primary',
    hollow: true,
  },

  {
    content: '进行重新登录流程',
    type: 'large',
    hollow: true,
  },
  {
    content:
      '如果卡在这里不出二维码，可以尝试 1 分钟后刷新页面，再次点击登录。如果还不行请删除此账号重新添加',
    type: 'large',
    hollow: true,
  },
];
const activities = ref([] as typeof fullActivities);

const store = useStore();
const message = useMessage();
const dialog = useDialog();

const accountTypes = computed(() => {
  const result = [
    {
      label: 'QQ(内置客户端)',
      value: 15,
      disabled: store?.diceServers?.length > 0 && store?.diceServers?.[0]?.baseInfo?.containerMode,
    },
    {
      label: 'QQ(内置 gocq)',
      value: 16,
      disabled: store?.diceServers?.length > 0 && store?.diceServers?.[0]?.baseInfo?.containerMode,
    },
    { label: 'QQ(Milky)', value: 17 },
    { label: 'QQ(onebot11 正向 WS)', value: 6 },
    { label: 'QQ(onebot11 反向 WS)', value: 11 },
    { label: 'Discord', value: 1 },
    { label: 'KOOK(开黑啦)', value: 2 },
    { label: 'Telegram', value: 3 },
  ];
  return result;
});

const curCaptchaIdSet = ref(''); // 当前设置了 ticket 的 id

const isRecentLogin = ref(false);
const duringRelogin = ref(false);
const dialogFormVisible = ref(false);
const dialogSlideVisible = ref(false);
const formLabelWidth = '120px';
const isTestMode = ref(false);

const slideIframe = ref(null);
const slideLink = ref('');
const slideBottomShow = ref(false);

const curConn = ref({} as DiceConnection);
const curConnId = ref('');

const signInfoLoaded = ref(false);
const signInfos = ref([] as SignInfo[]);
const signVerWarningText = ref('');
const signServerWarningText = ref('');

const setRecentLogin = () => {
  isRecentLogin.value = true;
  // 无用
  // curConn.value.adapter.inPackGoCqHttpRunning = false;
  // curConn.value.adapter.inPackGoCqHttpLoginDeviceLockUrl = '';
  setTimeout(() => {
    isRecentLogin.value = false;
  }, 3000);
};

const openSocks = async () => {
  const ret = await postToolOnebot();
  if (ret.ok) {
    dialog.success({
      title: '开启辅助工具',
      positiveText: '确定',
      content: () => (
        <>
          <n-text tag='p'>将在服务器上开启临时 socks5 服务，端口 13325</n-text>
          <n-text tag='p'>默认持续时长为 20 分钟</n-text>
          <n-text tag='p'>
            可能的公网 IP：
            <n-text type='success'>{ret.ip}</n-text>
          </n-text>
          <n-text tag='p' type='warning'>
            注：ip 不一定对仅供参考
          </n-text>
          <n-text tag='p'>请于服务器管理面板放行 13325 端口，协议 TCP</n-text>
          <n-text tag='p'>
            如果为 Windows Server 系统，请再额外关闭系统防火墙或设置放行规则。
          </n-text>
        </>
      ),
    });
  } else {
    dialog.error({
      title: '开启辅助工具',
      positiveText: '确定',
      content: () => (
        <>
          <n-text tag='p'>启动服务失败，或已经启动</n-text>
          <n-text tag='p'>
            报错信息：
            <n-text type='error'>{ret.errText}</n-text>
          </n-text>
          <n-text tag='p'>
            可能的公网 IP：
            <n-text type='success'>{ret.ip}</n-text>
          </n-text>
          <n-text tag='p' type='warning'>
            注：ip 不一定对仅供参考
          </n-text>
          <n-text tag='p'>请于服务器管理面板放行 13325 端口，协议 TCP</n-text>
          <n-text tag='p'>
            如果为 Windows Server 系统，请再额外关闭系统防火墙或设置放行规则。
          </n-text>
        </>
      ),
    });
  }
};

const goStepTwo = async () => {
  form.step = 2;
  curConnId.value = '';
  setRecentLogin();
  duringRelogin.value = false;

  store
    .addImConnection(form as any)
    .then(conn => {
      if ((conn as any).testMode) {
        isTestMode.value = true;
      } else {
        curConnId.value = conn.id;
      }
    })
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    .catch(e => {
      dialog.error('似乎已经添加了这个账号！', { title: '添加失败' });
      formClose();
    });
  if (form.accountType > 0) {
    dialogFormVisible.value = false;
    form.account = '';
    form.step = 1;
    return;
  }
  activities.value = [];
  await sleep(500);
  activities.value.push(fullActivities[0]);
  await sleep(1000);
  activities.value.push(fullActivities[1]);
  await sleep(1000);
  activities.value.push(fullActivities[2]);
};

const formClose = async () => {
  curConnId.value = '';
  dialogFormVisible.value = false;
  form.step = 1;
  form.isEnd = false;
};

const setEnable = async (i: DiceConnection, val: boolean) => {
  const ret = await postConnectSetEnable(i.id, val);
  i.enable = ret.enable;
  curCaptchaIdSet.value = '';
  message.success('状态修改完成');
  if (val) {
    setRecentLogin();
    // 若是启用骰子，走登录流程
    curConnId.value = ''; // 先改掉这个，以免和当前连接一致，导致被瞬间重置
    await nextTick(() => {
      curConnId.value = i.id;
    });
    // store.gocqhttpReloginImConnection(i).then(theConn => {
    //   curConnId.value = i.id;
    // })

    // 重复登录时，也打开这个窗口
    activities.value = [];
    dialogFormVisible.value = true;

    if (i.adapter.useInPackGoCqhttp) {
      form.step = 2;
      activities.value.push(fullActivities[4]);
      activities.value.push(fullActivities[5]);
      activities.value.push(fullActivities[2]);
    } else {
      form.step = 4;
      form.isEnd = true;
    }
  }
};

const askSetEnable = async (i: DiceConnection, val: boolean) => {
  dialog.warning({
    title: '警告',
    content: '确认修改此账号的在线状态吗',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await setEnable(i, val);
    },
  });
};

const askGocqhttpReLogin = async (i: DiceConnection) => {
  duringRelogin.value = false;
  dialog.warning({
    title: '警告',
    content: '重新登录吗？有可能要再次扫描二维码',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      gocqhttpReLogin(i);
    },
  });
};

const doGocqExport = async (i: DiceConnection) => {
  duringRelogin.value = false;
  dialog.warning({
    title: '提示',
    content: '即将下载 gocq 配置，是否继续？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      // http://localhost:3211/sd-api/im_connections/gocq_config_download.zip?id=10f576a4-5237-43f6-9086-269a9f9aace5&token=J4JAofWluYsc0YTgUtDuw3eBnVbZyW%232gTG0agA%40aAVRRIFmrTT0w4tEMbVxGdXn%3A0000000063a8664f
      location.href = `${urlBase}/sd-api/im_connections/gocq_config_download.zip?token=${encodeURIComponent(store.token)}&id=${encodeURIComponent(i.id)}`;
    },
  });
};

const gocqhttpReLogin = async (i: DiceConnection) => {
  curCaptchaIdSet.value = '';
  setRecentLogin();
  duringRelogin.value = true;
  curConnId.value = ''; // 先改掉这个，以免和当前连接一致，导致被瞬间重置
  if (curConn.value && curConn.value.adapter) {
    curConn.value.adapter.loginState = goCqHttpStateCode.Init;
  }
  postGoCqHttpRelogin(i.id)
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    .then(theConn => {
      curConnId.value = i.id;
    })
    .finally(() => {
      form.isEnd = true;
    });
  // 重复登录时，也打开这个窗口
  activities.value = [];
  dialogFormVisible.value = true;

  if (i.adapter.useInPackGoCqhttp) {
    form.step = 2;
    activities.value.push(fullActivities[4]);
    activities.value.push(fullActivities[5]);
    activities.value.push(fullActivities[2]);
  } else {
    form.step = 4;
  }
};

const signServerVersionChange = () => {
  switch (form.signServerVersion) {
    case '自定义':
      form.signServerName = '';
      signVerWarningText.value = '';
      signServerWarningText.value = '';
      break;
    case '':
      signInfos.value.forEach(info => {
        if (signInfos.value.length > 0) {
          form.signServerVersion = signInfos.value[0].version;
        }
        if (info.selected) {
          form.signServerVersion = info.version;
          signServerVersionChange();
        }
      });
      break;
    default:
      signInfos.value.forEach(info => {
        if (form.signServerVersion === info.version) {
          if (info.servers.length > 0) {
            form.signServerName = info.servers[0].name;
          }
          info.servers.forEach(server => {
            if (server.selected) {
              form.signServerName = server.name;
            }
          });
        }
      });
      break;
  }
  signServerServerChange();
};

const signServerServerChange = () => {
  signInfos.value.forEach(info => {
    if (info.version === form.signServerVersion) {
      signVerWarningText.value = info.note;
    }
    info.servers.forEach(server => {
      if (server.name === form.signServerName) {
        signServerWarningText.value = server.note;
      }
    });
  });
};

const getSignInfo = async () => {
  form.account = '';
  form.signServerVersion = '';
  form.signServerName = '';
  try {
    const res = await getLagrangeSignInfo();
    signInfoLoaded.value = res.result;
    // 理论上不会有 false 出现所以前端不进行通知，若出现则为后端代码问题
    if (res.result !== false) {
      signInfos.value = res.data;
      signServerVersionChange();
    }
  } catch {
    signInfoLoaded.value = false;
  }
};

const supportedQQVersions = ref<string[]>([]);

const form = reactive({
  accountType: 15,
  step: 1,
  isEnd: false,
  account: '',
  nickname: '',
  password: '',
  protocol: 1,
  appVersion: '',
  implementation: '',
  id: '',
  token: '',
  botToken: '',
  appToken: '',
  proxyURL: '',
  reverseProxyUrl: '',
  reverseProxyCDNUrl: '',
  url: '',
  clientID: '',
  robotCode: '',
  ignoreFriendRequest: false,
  extraArgs: '',
  endpoint: null as any as DiceConnection,

  relWorkDir: '',
  accessToken: '',
  connectUrl: '',

  host: '',
  port: '',

  appID: undefined,
  appSecret: '',
  onlyQQGuild: true,

  useSignServer: false,
  signServerConfig: {
    signServers: [
      {
        url: '',
        key: '',
        authorization: '',
      },
    ],
    ruleChangeSignServer: 1,
    maxCheckCount: 0,
    signServerTimeout: 60,
    autoRegister: false,
    autoRefreshToken: false,
    refreshInterval: 40,
  },
  signServerName: '',
  signServerKey: '',
  signServerVersion: '',

  reverseAddr: ':4001',
  platform: 'QQ',
  // milky
  wsGateway: '',
  restGateway: '',
});

export type addImConnectionForm = typeof form;

// 添加一个新账号
const addOne = () => {
  dialogFormVisible.value = true;
  form.protocol = 6;
  form.implementation = 'gocq';

  // 从后端读取 signInfo，动态填充到页面
  getSignInfo();
};

let timerId: number;

onBeforeMount(async () => {
  await store.getImConnections();
  for (const i of store.curDice.conns || []) {
    delete store.curDice.qrcodes[i.id];
  }

  // form.accountType 默认账号类型，在 android 与 mac 系统中，默认账号类型为内置 gocq，其余系统为内置客户端
  if (store.diceServers.length > 0) {
    if (
      store.diceServers[0].baseInfo.OS === 'android' ||
      store.diceServers[0].baseInfo.OS === 'darwin'
    ) {
      form.accountType = 16;
    }
    if (store.diceServers[0].baseInfo.containerMode) {
      form.accountType = 6;
    }
  }

  timerId = setInterval(async () => {
    console.log('refresh');
    await store.getImConnections();

    for (const i of store.curDice.conns || []) {
      // 获取二维码
      if (i.adapter?.loginState === goCqHttpStateCode.InLoginQrCode) {
        store.curDice.qrcodes[i.id] = (await postConnectionQrcode(i.id)).img;
      }

      if (i.id === curConnId.value) {
        curConn.value = i;

        // 登录失败
        if (i.state !== 1 && i.adapter?.loginState === goCqHttpStateCode.LoginFailed) {
          form.isEnd = true;
        }

        // 登录成功
        if (i.state === 1 && i.adapter?.loginState === goCqHttpStateCode.LoginSuccessed) {
          activities.value.push(fullActivities[3]);
          await sleep(1000);
          form.step = 3;
          form.isEnd = true;
        }

        break;
      }
    }
  }, 3000) as any;
});

onBeforeUnmount(() => {
  clearInterval(timerId);
});

const doRemove = async (i: DiceConnection) => {
  dialog.warning({
    title: '删除账号',
    content: '删除此项帐号，确定吗？（注：删除账号不会影响人物卡和 logs 等数据）',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await postConnectionDel(i.id);
      await store.getImConnections();
      message.success('删除成功！');
    },
  });
};
</script>

<style scoped>
.btn-add {
  width: 3rem !important;
  height: 3rem !important;
  font-size: 2rem;
  font-weight: bold;
}
</style>
