  <template>
    <div id="form-popup" v-if="formDisplay" class="container mt-3">
        <form id="submitForm">
            <div class="row mb-4">
                <div class="col-sm-12">
                    <h3 class="text-center">會議室名稱：{{ roomInfo.roomName }}</h3>
                </div>
            </div>
            <div class="row mb-2">
                <ItemPeriod 
                    period-name="早上："
                    :reservator-list="roomInfo.reservatorList.slice(0, 6)"
                    :info-progress-width="250"
                    :margin-left='170'
                    @update-show-reservator="updateShowReservator"
                    @showDiv="$emit('showDiv', $event)"
                    @hideDiv="$emit('hideDiv')"
                />
                <ItemPeriod 
                    period-name="下午："
                    :reservator-list="roomInfo.reservatorList.slice(6)"
                    :info-progress-width="500"
                    :margin-left='170'
                    @update-show-reservator="updateShowReservator"
                    @showDiv="$emit('showDiv', $event)"
                    @hideDiv="$emit('hideDiv')"
                />
            </div>
            <div class="row mb-2">
                <label for="name" class="col-sm-2 col-form-label">會議標題：</label>
                <div class="col-sm-9">
                    <input type="text" id="name" name="name" class="form-control" v-model="formInfo.name" required>
                </div>
            </div>

            <div class="row mb-2">
                <label for="email" class="col-sm-2 col-form-label">參與人郵件：</label>
                <div class="col-sm-9">
                    <input type="text" id="email" name="email" class="form-control" v-model="formInfo.email" required>
                </div>
            </div>

            <div class="row mb-2">
                <label class="col-sm-2 col-form-label">開始時間：</label>
                <div class="col-sm-9">
                    <select class="form-select" v-model="formInfo.start_time" required>
                        <option :value="item" v-for="item in time_period" :key="item">{{ item }}</option>
                    </select>
                </div>
            </div>

            <div class="row mb-2">
                <label class="col-sm-2 col-form-label">結束時間：</label>
                <div class="col-sm-9">
                    <select class="form-select" v-model="formInfo.end_time" required>
                        <option :value="item" v-for="item in time_period" :key="item">{{ item }}</option>
                    </select>
                </div>
            </div>

            <div class="row mb-2">
                <label for="content" class="col-sm-2 col-form-label">會議內容：</label>
                <div class="col-sm-9">
                    <textarea id="content" name="content" class="form-control" maxlength="200" v-model="formInfo.content"></textarea>
                </div>
            </div>

            <div class="row mb-2">
                <label for="file" class="col-sm-2 col-form-label">上傳文件：</label>
                <div class="col-sm-9">
                    <input type="file" id="file" name="file" class="form-control">
                </div>
            </div>

            <div class="row mb-2">
                <div class="col-sm-12 d-flex justify-content-center">
                    <button type="submit" class="btn btn-primary margin-right-2cm" @click.prevent="submitForm">預約</button>
                    <button type="button" class="btn btn-secondary" @click="closeForm">取消</button>
                </div>
            </div>
        </form>
    </div>
  </template>
  
  <script>
  import ItemPeriod from './ItemPeriod.vue';

  export default {
    name: 'ReserveForm',
    emits: ['close-form', 'showDiv', 'hideDiv'],
    props: ['roomInfo', 'formDisplay', 'userName'],
    data() {
      return {
        formInfo: {
          name: this.userName,
          email: 'example@gmail.com',
          start_time: '10:00',
          end_time: '12:00',
          content: 'test content',
          file: ''
        },
        time_period: [],
        showReservator: '',
      };
    },
    methods: {
      submitForm() {
        console.log(this.formInfo);
        console.log(this.roomInfo);
        
        this.closeForm();
      },
      updateShowReservator(value) {
        this.showReservator = value;
      },
      closeForm() {
        this.$emit('close-form');
      }
    },
    components: {
        ItemPeriod
    },
    mounted() {
      this.time_period = this.$names;
    }
  }
  </script>
  