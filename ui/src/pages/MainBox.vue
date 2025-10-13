<template>
  <n-layout id="root" class="mx-auto my-0 flex h-screen flex-col" position="absolute">
    <n-layout-header bordered class="nav flex h-16 flex-none justify-between bg-inherit">
      <n-flex :size="0" align="center" style="height: 60px">
        <div class="menu-button-wrapper mx-2 flex justify-center">
          <n-button size="large" text @click="drawerMenu = true">
            <n-icon size="1.5rem">
              <i-carbon-menu />
            </n-icon>
          </n-button>
        </div>

        <n-flex :size="0" :v-show="store.canAccess" align="flex-start" vertical>
          <n-flex align="center" size="small">
            <n-text style="font-size: 1.2rem; cursor: pointer" @click="enableAdvancedConfig">
              {{ store.diceServers[0].baseInfo.appName ?? 'SealDice-CE' }}
            </n-text>
            <n-tooltip
              v-if="store.diceServers.length > 0 && store.diceServers[0].baseInfo.containerMode"
              class="flex items-center">
              <template #trigger>
                <n-button text>
                  <n-icon>
                    <i-carbon-container-software />
                  </n-icon>
                </n-button>
              </template>
              当前以容器模式启动，部分功能受到限制。
            </n-tooltip>
          </n-flex>
          <n-text v-if="store.diceServers.length > 0" style="font-size: 0.7rem">
            {{ store.diceServers[0].baseInfo.OS }} - {{ store.diceServers[0].baseInfo.arch }}
          </n-text>
        </n-flex>
      </n-flex>

      <n-flex v-show="store.canAccess" align="center" size="medium">
        <div style="cursor: pointer">
          <n-badge :show="!newsChecked" value="new">
            <n-button text @click="dialogFeed = true" :type="newsChecked ? 'default' : 'error'">
              <n-icon size="1.75rem">
                <i-carbon-notification />
              </n-icon>
            </n-button>
          </n-badge>
        </div>

        <div style="display: flex; flex-direction: column; align-items: center">
          <div style="display: flex; align-items: center">
            <n-tag
              :bordered="false"
              :type="store.curDice.baseInfo.appChannel === 'stable' ? 'success' : 'default'"
              size="small"
              style="margin-right: 0.3rem">
              {{ store.curDice.baseInfo.appChannel === 'stable' ? '正式版' : '测试版' }}
            </n-tag>
            <n-tooltip placement="bottom">
              <template #trigger>
                <n-text class="text-base">
                  {{ store.curDice.baseInfo.versionSimple }}
                </n-text>
              </template>
              {{ store.curDice.baseInfo.version }}
            </n-tooltip>
          </div>
          <div v-if="store.curDice.baseInfo.versionCode < store.curDice.baseInfo.versionNewCode">
            🆕{{ store.curDice.baseInfo.versionNew }}
          </div>
        </div>

        <n-button quaternary circle @click="toggleDark()">
          <template #icon>
            <n-icon>
              <i-carbon-asleep v-if="isDark" />
              <i-carbon-light v-else />
            </n-icon>
          </template>
        </n-button>
      </n-flex>
    </n-layout-header>

    <n-layout class="mt-16 flex grow" position="absolute" has-sider>
      <n-layout-sider
        collapse-mode="width"
        :collapsed-width="64"
        :width="240"
        :collapsed="collapsedMenu"
        show-trigger
        @collapse="collapsedMenu = true"
        @expand="collapsedMenu = false"
        bordered
        :native-scrollbar="false"
        class="menu no-scrollbar flex-none bg-inherit">
        <Menu
          v-model:advancedConfigCounter="advancedConfigCounter"
          :collapsed="collapsedMenu"
          type="dark" />
      </n-layout-sider>

      <n-layout-content
        class="h-auto flex-1 overflow-y-auto text-left"
        :native-scrollbar="false"
        content-class="h-full"
        embedded>
        <n-spin
          ref="rightbox"
          :show="loading"
          class="main-container size-full"
          content-class="h-full">
          <router-view
            v-if="!loading"
            @update:advanced-settings-show="(show: boolean) => refreshAdvancedSettings(show)" />
        </n-spin>
      </n-layout-content>
    </n-layout>
  </n-layout>

  <n-drawer v-model:show="drawerMenu" class="drawer-menu" default-width="50%" placement="left">
    <n-drawer-content body-content-style="padding: 0;" :native-scrollbar="false">
      <template #header>
        <n-flex size="small">
          <n-flex :size="0" :v-show="store.canAccess" align="flex-start" vertical>
            <span class="cursor-pointer text-base" @click="enableAdvancedConfig">
              {{ store.diceServers[0].baseInfo.appName }}?? 'SealDice-CE' }}
            </span>
            <span v-if="store.diceServers.length > 0" class="text-xs">
              {{ store.diceServers[0].baseInfo.OS }} - {{ store.diceServers[0].baseInfo.arch }}
            </span>
          </n-flex>
          <n-tag
            :bordered="false"
            :type="store.curDice.baseInfo.appChannel === 'stable' ? 'success' : 'default'"
            size="small">
            {{ store.curDice.baseInfo.appChannel === 'stable' ? '正式版' : '测试版' }}
          </n-tag>
        </n-flex>
      </template>

      <Menu v-model:advancedConfigCounter="advancedConfigCounter" />
    </n-drawer-content>
  </n-drawer>

  <n-modal
    v-model:show="showDialog"
    preset="dialog"
    :closable="false"
    :mask-closable="false"
    :close-on-esc="false"
    title="输入密码解锁"
    positive-text="确认"
    @positive-click="doUnlock"
    @keyup.enter="doUnlock">
    <n-input v-model:value="password" type="password" />
  </n-modal>

  <n-modal
    v-model:show="dialogLostConnectionVisible"
    :closable="false"
    preset="card"
    style="width: 30%"
    title="主程序离线"
    :mask-closable="false"
    transform-origin="center">
    <n-text>与主程序的连接断开，请耐心等待连接恢复。如果失去响应过久，请登录服务器处理。</n-text>
  </n-modal>

  <n-modal
    v-model:show="dialogCheckPassword"
    :closable="false"
    preset="dialog"
    style="width: 30%"
    title="欢迎使用海豹核心"
    :mask-closable="false"
    transform-origin="center">
    <n-text>
      如果您的服务开启在公网，为了保证您的安全性，请前往
      <strong>「综合设置」>「基本设置」</strong> 界面，设置
      <strong>UI 界面密码</strong>。或切换为只有本机可访问。<br />
    </n-text>
    <n-text type="warning" class="mt-4 block">
      如果您不了解上面在说什么，请务必设置一个密码！
    </n-text>

    <template #action>
      <n-button type="primary" :disabled="!canSkip" @click="dialogCheckPassword = false">
        我已知晓！
        <template v-if="!canSkip">
          （<n-countdown
            :duration="5 * 1000"
            :render="renderCheckPasswordCountDown"
            @finish="canSkip = true" />
          秒后可点击）
        </template>
      </n-button>
    </template>
  </n-modal>

  <n-modal
    v-model:show="dialogFeed"
    :closable="false"
    :mask-closable="false"
    class="w-3/4"
    preset="card"
    title="海豹新闻">
    <template #header-extra>
      <n-button type="primary" @click="checkNews">
        <template #icon>
          <n-icon>
            <i-carbon-checkmark />
          </n-icon>
        </template>
        确认已读
      </n-button>
    </template>

    <div style="text-align: left" v-html="newsData"></div>
  </n-modal>
</template>

<script setup lang="tsx">
import { useStore } from '~/store';
import { useMessage, useNotification } from 'naive-ui';
import { passwordHash } from '~/utils';
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
import relativeTime from 'dayjs/plugin/relativeTime';
import { getNewUtils, postUtilsCheckNews } from '~/api/v1/utils';
import { checkSecurity } from '~/api/v1/others';
import { breakpointsTailwind } from '@vueuse/core';

const isDark = useDark({ disableTransition: false });
const toggleDark = useToggle(isDark);
const message = useMessage();
const notification = useNotification();

const breakpoints = useBreakpoints(breakpointsTailwind);
const isDesktop = breakpoints.greater('md');

dayjs.locale('zh-cn');
dayjs.extend(relativeTime);

const loading = useStorage('router-view-loading', true);

const collapsedMenu = ref<boolean>(!isDesktop.value);

const store = useStore();
const password = ref('');

const dialogFeed = ref(false);

const newsData = ref(`<div>暂无内容</div>`);
const newsChecked = ref(true);
const newsMark = ref('');
const checkNews = async () => {
  const ret = await postUtilsCheckNews(newsMark.value);
  if (ret?.result) {
    message.success('已阅读最新的海豹新闻');
  } else {
    message.error('阅读海豹新闻失败');
  }
  dialogFeed.value = false;
  await updateNews();
};
const updateNews = async () => {
  const newsInfo = await getNewUtils();
  if (newsInfo.result) {
    newsData.value = newsInfo.news;
    newsChecked.value = newsInfo.checked;
    newsMark.value = newsInfo.newsMark;
  } else {
    message.error(newsInfo?.err ?? '获取海豹新闻失败');
  }
};

const showDialog = computed(() => {
  return !store.canAccess;
});

const dialogLostConnectionVisible = ref(false);

const doUnlock = async () => {
  const hash = await passwordHash(store.salt, password.value);
  await store.signIn(hash);
  if (store.canAccess) {
    password.value = '';
    notification['success']({
      title: '登录成功！',
      content: '欢迎回来，请开始使用',
      duration: 3000,
      closable: false,
    });
    await checkPassword();
    window.location.reload();
  } else {
    notification['error']({
      title: '登录失败',
      content: '密码错误',
      duration: 3000,
      closable: false,
    });
    password.value = '';
  }
};

const dialogCheckPassword = ref(false);
const canSkip = ref<boolean>(false);
const checkPassword = async () => {
  if (!(await checkSecurity()).isOk) {
    dialogCheckPassword.value = true;
    canSkip.value = false;
  }
};
const renderCheckPasswordCountDown = ({ seconds }) => <span>{seconds}</span>;

onBeforeMount(async () => {
  store.getBaseInfo();
  store.getCustomText();

  if (store.canAccess) {
    checkPassword();
  }

  timerId = setInterval(async () => {
    try {
      await store.getBaseInfo();
      if (dialogLostConnectionVisible.value) {
        dialogLostConnectionVisible.value = false;
      }
    } catch (e: any) {
      if (!e.response) {
        // 此时是连接不上，404
        // e.response.status 有可能为 403
        dialogLostConnectionVisible.value = true;
      }
    }
  }, 5000) as any;

  await updateNews();

  const conf = await store.diceAdvancedConfigGet();
  if (conf.show) {
    advancedConfigCounter.value = 8;
  }
});

onBeforeUnmount(() => {
  clearInterval(timerId);
});

let timerId: number;

const rightbox = ref(null);

const drawerMenu = ref<boolean>(false);

const advancedConfigCounter = ref<number>(0);
const enableAdvancedConfig = async () => {
  advancedConfigCounter.value++;
  const counter = advancedConfigCounter.value;
  if (counter > 8) {
    ElMessage.info('高级设置页已经开启');
    await router.push({ path: '/misc/advanced-setting' });
    return;
  } else if (counter === 8) {
    const conf = await store.diceAdvancedConfigGet();
    conf.show = true;
    conf.enable = true;
    await store.diceAdvancedConfigSet(conf);
    await router.push({ path: '/misc/advanced-setting' });
    ElMessage.success('已开启高级设置页');
  } else if (counter > 2) {
    ElMessage.info('再按 ' + (8 - counter) + ' 次开启高级设置页');
  }
};

const router = useRouter();
const refreshAdvancedSettings = async (show: boolean) => {
  if (!show) {
    advancedConfigCounter.value = 0;
    await router.push({ path: '/log', replace: true });
    ElMessage.success('已关闭高级设置页');
  }
};
</script>

<style lang="css">
html,
body {
  height: 100%;
}

.main-container {
  padding: 2rem;
  box-sizing: border-box;
  min-height: 100%;
}

@media screen and (max-width: 640px) {
  .nav {
    padding: 0 0.5rem 0 0;
  }

  .menu {
    display: none;
  }

  .main-container {
    padding: 1rem;
  }
}

@media screen and (min-width: 640px) {
  .nav {
    padding: 0 1rem 0 1.5rem;
  }

  .menu-button-wrapper {
    display: none;
  }
}

#app {
  font-family:
    'PingFang SC', 'Helvetica Neue', 'Hiragino Sans GB', 'Segoe UI', 'Microsoft YaHei', '微软雅黑',
    sans-serif;
}
</style>
