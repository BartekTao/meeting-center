<template>
  <comm-with-gql ref="commWithGql"></comm-with-gql>
  <div id="form-popup" v-if="formDisplay" class="container mt-3">
    <form id="submitForm">
      <div class="row mb-2">
        <label for="name" class="col-sm-2 col-form-label">會議名稱：</label>
        <div class="col-sm-9">
          <input type="text" id="name" name="name" class="form-control" v-model="localFormInfo.name" required>
        </div>
      </div>

      <div class="row mb-2">
        <label for="image_url" class="col-sm-2 col-form-label">圖片網址：</label>
        <div class="col-sm-9">
          <input type="text" id="image_url" name="image_url" class="form-control" v-model="image_url" required>
        </div>
      </div>

      <!-- <div class="row mb-2">
        <label for="people_limit" class="col-sm-2 col-form-label">人數限制：

        </label>
        <div class="col-sm-9">
          <input type="number" id="people_limit" name="people_limit" class="form-control" v-model="localFormInfo.capacity" required>
        </div>
      </div> -->

      <div class="row mb-2">
        <label for="can_eat" class="col-sm-2 col-form-label">可否進食：</label>
        <div class="col-sm-9">
          <input type="checkbox" id="can_eat" name="can_eat" class="form-check-input mt-2" v-model="canFood">
        </div>
      </div>

      <!-- <div class="row mb-2">
        <label for="can_drink" class="col-sm-2 col-form-label">可否喝水：</label>
        <div class="col-sm-9">
          <input type="checkbox" id="can_drink" name="can_drink" class="form-check-input mt-2" v-model="canDrink">
        </div>
      </div> -->

      <div class="row mb-2">
        <label for="has_big_table" class="col-sm-2 col-form-label">有白板：</label>
        <div class="col-sm-9">
          <input type="checkbox" id="has_big_table" name="has_big_table" class="form-check-input mt-2" v-model="hasBigTable">
        </div>
      </div>

      <div class="row mb-2">
        <label for="has_projector" class="col-sm-2 col-form-label">有投影機：</label>
        <div class="col-sm-9">
          <input type="checkbox" id="has_projector" name="has_projector" class="form-check-input mt-2" v-model="hasProjector">
        </div>
      </div>
      <div class="row mb-2">
        <div class="col-sm-12 d-flex justify-content-center">
          <button type="submit" class="btn btn-primary margin-right-2cm" @click.prevent="submitForm">編輯</button>
          <button type="button" class="btn btn-secondary" @click="closeForm">取消</button>
        </div>
      </div>
    </form>
  </div>
</template>

<script>
import CommWithGql from '@/components/CommWithGql.vue'

export default {
    components: {
      CommWithGql,
    },
  emits: ['close-form', 'updateInfo', 'update-all-rooms'],
  props: ['roomInfo', 'formDisplay'],
  data() {
    return {
      image_url: require('../assets/images/listing-01.jpg'),
    };
  },
  computed: {
    localFormInfo() {
      return JSON.parse(JSON.stringify(this.roomInfo));
    },
    canDrink: {
      get() {
        return !this.localFormInfo.rules.includes('NO_DRINK');
      },
      set(value) {
        const index = this.localFormInfo.rules.indexOf('NO_DRINK');
        if (value && index !== -1) {
          this.localFormInfo.rules.splice(index, 1);
        } else if (!value && index === -1) {
          this.localFormInfo.rules.push('NO_DRINK');
        }
      }
    },
    canFood: {
      get() {
        return !this.localFormInfo.rules.includes('NO_FOOD');
      },
      set(value) {
        const index = this.localFormInfo.rules.indexOf('NO_FOOD');
        if (value && index !== -1) {
          this.localFormInfo.rules.splice(index, 1);
        } else if (!value && index === -1) {
          this.localFormInfo.rules.push('NO_FOOD');
        }
      }
    },
    hasBigTable: {
      get() {
        return this.localFormInfo.equipments.includes('TABLE');
      },
      set(value) {
        const index = this.localFormInfo.equipments.indexOf('TABLE');
        if (!value && index !== -1) {
          this.localFormInfo.equipments.splice(index, 1);
        } else if (value && index === -1) {
          this.localFormInfo.equipments.push('TABLE');
        }
      }
    },
    hasProjector: {
      get() {
        return this.localFormInfo.equipments.includes('PROJECTOR');
      },
      set(value) {
        const index = this.localFormInfo.equipments.indexOf('PROJECTOR');
        if (!value && index !== -1) {
          this.localFormInfo.equipments.splice(index, 1);
        } else if (value && index === -1) {
          this.localFormInfo.equipments.push('PROJECTOR');
        }
      }
    }
  },
  methods: {
    submitForm() {
      this.$refs.commWithGql.createRoom(this.localFormInfo)
      this.$emit('update-all-rooms');
      this.closeForm();
    },
    closeForm() {
      this.$emit('close-form');
    },
    updateValue(field, value) {
      this.$emit('update-info', { field, value });
    },
  }
};
</script>
