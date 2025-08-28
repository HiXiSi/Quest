<template>
  <div class="model-viewer">
    <div ref="viewerContainer" class="viewer-container"></div>
    <div class="viewer-controls">
      <el-button-group>
        <el-button @click="resetCamera" size="small">
          <el-icon><Refresh /></el-icon>
          重置视角
        </el-button>
        <el-button @click="toggleWireframe" size="small">
          <el-icon><Grid /></el-icon>
          {{ showWireframe ? '实体' : '线框' }}
        </el-button>
        <el-button @click="toggleAnimation" size="small" v-if="hasAnimations">
          <el-icon><VideoPlay /></el-icon>
          {{ animationPlaying ? '暂停' : '播放' }}
        </el-button>
      </el-button-group>
    </div>
    <div v-if="loading" class="loading-overlay" v-loading="true" element-loading-text="正在加载模型...">
    </div>
    <div v-if="error" class="error-overlay">
      <el-icon size="48" color="#f56c6c"><WarningFilled /></el-icon>
      <div class="error-text">{{ error }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Grid, VideoPlay, WarningFilled } from '@element-plus/icons-vue'
import * as THREE from 'three'
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js'
import { FBXLoader } from 'three/examples/jsm/loaders/FBXLoader.js'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'

const props = defineProps({
  modelUrl: {
    type: String,
    required: true
  },
  fileExtension: {
    type: String,
    required: true
  }
})

const viewerContainer = ref()
const loading = ref(true)
const error = ref('')
const showWireframe = ref(false)
const hasAnimations = ref(false)
const animationPlaying = ref(false)

let scene, camera, renderer, controls, model, mixer, clock
let animationId

const initThreeJS = () => {
  // 创建场景
  scene = new THREE.Scene()
  scene.background = new THREE.Color(0xf0f0f0)

  // 创建相机
  camera = new THREE.PerspectiveCamera(75, 1, 0.1, 1000)
  camera.position.set(5, 5, 5)

  // 创建渲染器
  renderer = new THREE.WebGLRenderer({ antialias: true })
  renderer.setSize(400, 300)
  renderer.shadowMap.enabled = true
  renderer.shadowMap.type = THREE.PCFSoftShadowMap
  
  viewerContainer.value.appendChild(renderer.domElement)

  // 创建控制器
  controls = new OrbitControls(camera, renderer.domElement)
  controls.enableDamping = true
  controls.dampingFactor = 0.1

  // 添加光源
  const ambientLight = new THREE.AmbientLight(0xf0f0f0, 5)
  scene.add(ambientLight)

  const directionalLight = new THREE.DirectionalLight(0xffffff, 8)
  directionalLight.position.set(100, 100, 50)
  directionalLight.castShadow = true
  scene.add(directionalLight)

  // 添加网格地面
  const gridHelper = new THREE.GridHelper(10, 10)
  gridHelper.material.opacity = 0.3
  gridHelper.material.transparent = true
  scene.add(gridHelper)

  // 创建时钟
  clock = new THREE.Clock()

  // 开始渲染循环
  animate()
}

const loadModel = async () => {
  try {
    loading.value = true
    error.value = ''

    const extension = props.fileExtension.toLowerCase()
    
    if (['gltf', 'glb'].includes(extension)) {
      await loadGLTFModel()
    } else if (extension === 'fbx') {
      await loadFBXModel()
    } else {
      throw new Error(`不支持的文件格式: ${extension}`)
    }

    // 调整相机位置以适应模型
    adjustCameraToModel()
    loading.value = false
  } catch (err) {
    console.error('模型加载失败:', err)
    error.value = err.message || '模型加载失败'
    loading.value = false
  }
}

const loadGLTFModel = () => {
  return new Promise((resolve, reject) => {
    const loader = new GLTFLoader()
    loader.load(
      props.modelUrl,
      (gltf) => {
        model = gltf.scene
        scene.add(model)
        
        // 检查是否有动画
        if (gltf.animations && gltf.animations.length > 0) {
          hasAnimations.value = true
          mixer = new THREE.AnimationMixer(model)
          gltf.animations.forEach((clip) => {
            const action = mixer.clipAction(clip)
            action.play()
          })
          animationPlaying.value = true
        }
        
        resolve()
      },
      (progress) => {
        console.log('加载进度:', (progress.loaded / progress.total * 100) + '%')
      },
      (error) => {
        reject(error)
      }
    )
  })
}

const loadFBXModel = () => {
  return new Promise((resolve, reject) => {
    const loader = new FBXLoader()
    loader.load(
      props.modelUrl,
      (fbx) => {
        model = fbx
        scene.add(model)
        
        // 检查是否有动画
        if (fbx.animations && fbx.animations.length > 0) {
          hasAnimations.value = true
          mixer = new THREE.AnimationMixer(model)
          fbx.animations.forEach((clip) => {
            const action = mixer.clipAction(clip)
            action.play()
          })
          animationPlaying.value = true
        }
        
        resolve()
      },
      (progress) => {
        console.log('加载进度:', (progress.loaded / progress.total * 100) + '%')
      },
      (error) => {
        reject(error)
      }
    )
  })
}

const adjustCameraToModel = () => {
  if (!model) return

  const box = new THREE.Box3().setFromObject(model)
  const size = box.getSize(new THREE.Vector3())
  const center = box.getCenter(new THREE.Vector3())

  // 计算适当的相机距离
  const maxDim = Math.max(size.x, size.y, size.z)
  const distance = maxDim * 2

  camera.position.copy(center)
  camera.position.x += distance
  camera.position.y += distance
  camera.position.z += distance
  camera.lookAt(center)

  controls.target.copy(center)
  controls.update()
}

const animate = () => {
  animationId = requestAnimationFrame(animate)
  
  controls.update()
  
  if (mixer && animationPlaying.value) {
    mixer.update(clock.getDelta())
  }
  
  renderer.render(scene, camera)
}

const resetCamera = () => {
  adjustCameraToModel()
  ElMessage.success('视角已重置')
}

const toggleWireframe = () => {
  if (!model) return
  
  showWireframe.value = !showWireframe.value
  
  model.traverse((child) => {
    if (child.isMesh) {
      child.material.wireframe = showWireframe.value
    }
  })
}

const toggleAnimation = () => {
  if (!mixer) return
  
  animationPlaying.value = !animationPlaying.value
  
  // 使用正确的Three.js AnimationAction API
  mixer._actions.forEach(action => {
    if (animationPlaying.value) {
      // 恢复播放：设置paused为false并保持enabled为true
      action.paused = false
      action.enabled = true
      if (!action.isRunning()) {
        action.play()
      }
    } else {
      // 暂停播放：设置paused为true
      action.paused = true
    }
  })
}

const handleResize = () => {
  if (!renderer || !camera) return
  
  const width = viewerContainer.value.clientWidth
  const height = viewerContainer.value.clientHeight
  
  camera.aspect = width / height
  camera.updateProjectionMatrix()
  renderer.setSize(width, height)
}

// 监听容器大小变化
let resizeObserver
onMounted(() => {
  initThreeJS()
  loadModel()
  
  resizeObserver = new ResizeObserver(handleResize)
  resizeObserver.observe(viewerContainer.value)
})

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId)
  }
  
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
  
  if (renderer) {
    renderer.dispose()
  }
  
  if (viewerContainer.value && renderer) {
    viewerContainer.value.removeChild(renderer.domElement)
  }
})

// 监听模型URL变化
watch(() => props.modelUrl, () => {
  if (model) {
    scene.remove(model)
  }
  loadModel()
})
</script>

<style scoped>
.model-viewer {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 400px;
}

.viewer-container {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
}

.viewer-controls {
  position: absolute;
  top: 16px;
  right: 16px;
  z-index: 10;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 20;
}

.error-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 20;
}

.error-text {
  margin-top: 16px;
  font-size: 14px;
  color: #f56c6c;
  text-align: center;
}
</style>