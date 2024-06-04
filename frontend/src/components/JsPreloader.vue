
<template>
    <div id="js-preloader" class="js-preloader" :class="{ 'loaded': isLoaded }">
        <div class="preloader-inner">
        <span class="dot"></span>
        <div class="dots">
            <span></span>
            <span></span>
            <span></span>
        </div>
        </div>
    </div>
</template>

<script>
    export default {
    name: 'JsPreloader',
    data() {
        return {
        isLoaded: true
        };
    },
    // mounted() {
    //     window.addEventListener('load', () => {
    //     this.isLoaded = true;
    //         });
    //     }
    };
</script>
