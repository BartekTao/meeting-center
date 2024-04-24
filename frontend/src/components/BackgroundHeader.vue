<template>
  <header 
    class="header-area header-sticky slideInDown" 
    :class="{'background-header': shouldAddBackground}" 
    ref="header">
    <div class="container">
      <div class="row">
        <div class="col-12">
          <nav class="main-nav">
            <ul class="nav">
              <li><router-link to="/" class="nav-link">空間查詢</router-link></li>
              <li><router-link to="/reserved-page" class="nav-link">已預定空間</router-link></li>
              <li><router-link :to="{name: '後台管理'}" class="nav-link">後台管理</router-link></li>
              <li><div class="main-white-button"><a href="#"><i class="fa fa-plus"></i> 會議預約</a></div></li> 
            </ul>        
            <!-- <a class='menu-trigger'>
              <span>Menu</span>
            </a> -->
          </nav>
        </div>
      </div>
    </div>
  </header>
</template>
  
  
  
<script>
export default {
  name: 'BackgroundHeader',
  data() {
    return {
      shouldAddBackground: false,
      isMenuVisible: false  // 用於追踪導航菜單的顯示狀態
    };
  },
  methods: {
    checkScroll() {
      const scroll = window.scrollY;
      // 確保 $refs.header 在模板中正確設置
      const headerHeight = this.$refs.header.offsetHeight;

      // 更新應用條件類的邏輯以應對變更的模板結構
      this.shouldAddBackground = scroll >= headerHeight;
    },
    toggleMenu() {
      // 切換菜單的顯示狀態
      this.isMenuVisible = !this.isMenuVisible;
    }
  },
  mounted() {
    window.addEventListener('scroll', this.checkScroll);
    // 確保在元件卸載前移除事件監聽器
  },
  beforeUnmount() {
    window.removeEventListener('scroll', this.checkScroll);
  }
};
</script>


  