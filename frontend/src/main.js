import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'

const app = createApp(App)

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(createPinia())
app.use(router)
app.use(ElementPlus)

// 初始化应用
const initApp = async () => {
  // 动态导入userStore以确保Pinia已初始化
  const { useUserStore } = await import('@/stores/user')
  const userStore = useUserStore()
  
  // 尝试初始化用户状态
  try {
    await userStore.initUser()
  } catch (error) {
    // 初始化失败不阻塞应用启动
    console.log('用户状态初始化失败:', error.message)
  }
  
  app.mount('#app')
}

initApp()