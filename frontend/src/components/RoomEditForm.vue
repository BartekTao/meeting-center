<template>
  <comm-with-gql ref="commWithGql"></comm-with-gql>
  <div id="form-popup" v-if="formDisplay" class="container mt-3">
    <form id="submitForm">
      <div class="row mb-2">
        <label for="name" class="col-sm-2 col-form-label">會議名稱：</label>
        <div class="col-sm-9">
          <input type="text" id="name" name="name" class="form-control" v-model="localFormInfo.roomId" required>
        </div>
      </div>

      <div class="row mb-2">
        <label for="image_url" class="col-sm-2 col-form-label">圖片網址：</label>
        <div class="col-sm-9">
          <input type="text" id="image_url" name="image_url" class="form-control" v-model="image_url" required>
        </div>
      </div>

      <div class="row mb-2">
        <label for="people_limit" class="col-sm-2 col-form-label">人數限制：

        </label>
        <div class="col-sm-9">
          <input type="number" id="people_limit" name="people_limit" class="form-control" v-model="localFormInfo.capacity" required>
        </div>
      </div>

      <div class="row mb-2">
        <label for="can_eat" class="col-sm-2 col-form-label">可否進食：</label>
        <div class="col-sm-9">
          <input type="checkbox" id="can_eat" name="can_eat" class="form-check-input mt-2" v-model="canFood">
        </div>
      </div>

      <div class="row mb-2">
        <label for="can_drink" class="col-sm-2 col-form-label">可否喝水：</label>
        <div class="col-sm-9">
          <input type="checkbox" id="can_drink" name="can_drink" class="form-check-input mt-2" v-model="canDrink">
        </div>
      </div>

      <div class="row mb-2">
        <label for="has_big_table" class="col-sm-2 col-form-label">有大桌子：</label>
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
          <!-- <button type="submit" class="btn btn-primary margin-right-2cm" >編輯</button>
          <button type="button" class="btn btn-secondary" >取消</button> -->
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
      CommWithGql
    },
  emits: ['close-form', 'updateInfo'],
  props: ['formInfo', 'formDisplay'],
  data() {
    return {
      image_url: require('../assets/images/listing-01.jpg'),
      new_room: {
        name: '',
        image_url: '',
        people_limit: 0,
        can_eat: false,
        can_drink: false,
        has_big_table: false,
        has_projector: false}
    };
  },
  computed: {
    localFormInfo() {
      return JSON.parse(JSON.stringify(this.formInfo));
    },
    canDrink: {
      get() {
        return !this.localFormInfo.rules.includes('no drinks');
      },
      set(value) {
        const index = this.localFormInfo.rules.indexOf('no drinks');
        if (value && index !== -1) {
          this.localFormInfo.rules.splice(index, 1);
        } else if (!value && index === -1) {
          this.localFormInfo.rules.push('no drinks');
        }
      }
    },
    canFood: {
      get() {
        return !this.localFormInfo.rules.includes('no food');
      },
      set(value) {
        const index = this.localFormInfo.rules.indexOf('no food');
        if (value && index !== -1) {
          this.localFormInfo.rules.splice(index, 1);
        } else if (!value && index === -1) {
          this.localFormInfo.rules.push('no food');
        }
      }
    },
    hasBigTable: {
      get() {
        return this.localFormInfo.equipment.includes('big table');
      },
      set(value) {
        const index = this.localFormInfo.equipment.indexOf('big table');
        if (!value && index !== -1) {
          this.localFormInfo.equipment.splice(index, 1);
        } else if (value && index === -1) {
          this.localFormInfo.equipment.push('big table');
        }
      }
    },
    hasProjector: {
      get() {
        return this.localFormInfo.equipment.includes('projector');
      },
      set(value) {
        const index = this.localFormInfo.equipment.indexOf('projector');
        if (!value && index !== -1) {
          this.localFormInfo.equipment.splice(index, 1);
        } else if (value && index === -1) {
          this.localFormInfo.equipment.push('projector');
        }
      }
    }
  },
  methods: {
    submitForm() {
      console.log(this.localFormInfo)
      this.$refs.commWithGql.createRoom(this.localFormInfo);
      this.$refs.commWithGql.QueryAllRooms();
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
