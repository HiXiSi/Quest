import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

// 路由组件懒加载
const Login = () => import('@/views/Login.vue')
const Register = () => import('@/views/Register.vue')
const Layout = () => import('@/views/Layout.vue')
const Dashboard = () => import('@/views/Dashboard.vue')
const FileManager = () => import('@/views/FileManager.vue')
const FileUpload = () => import('@/views/FileUpload.vue')
const CategoryManager = () => import('@/views/CategoryManager.vue')
const TagManager = () => import('@/views/TagManager.vue')
const RecycleBin = () => import('@/views/RecycleBin.vue')
const UserManager = () => import('@/views/UserManager.vue')
const Profile = () => import('@/views/Profile.vue')

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: Layout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: Dashboard
      },
      {
        path: 'files',
        name: 'FileManager',
        component: FileManager
      },
      {
        path: 'upload',
        name: 'FileUpload',
        component: FileUpload
      },
      {
        path: 'categories',
        name: 'CategoryManager',
        component: CategoryManager
      },
      {
        path: 'tags',
        name: 'TagManager',
        component: TagManager
      },
      {
        path: 'recycle',
        name: 'RecycleBin',
        component: RecycleBin
      },
      {
        path: 'users',
        name: 'UserManager',
        component: UserManager,
        meta: { requiresAdmin: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: Profile
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
  } else if (to.meta.requiresAdmin && userStore.user?.role !== 'admin') {
    next('/')
  } else if ((to.name === 'Login' || to.name === 'Register') && userStore.isLoggedIn) {
    next('/')
  } else {
    next()
  }
})

export default router