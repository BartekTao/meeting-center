  <template>
    <comm-with-gql ref="commWithGql"></comm-with-gql>
  
    <div id="form-popup" v-if="formDisplay" class="container mt-3">
        <form id="submitForm">
            <div class="row mb-4">
                <div class="col-sm-12">
                    <h3 class="text-center">會議室名稱：{{ localFormInfo.name }}</h3>
                </div>
            </div>
            <div class="row mb-2">
                <ItemPeriod 
                    period-name="早上："
                    :reservator-list="schedulesList.slice(0, 6)"
                    :info-progress-width="250"
                    :margin-left='170'
                    @update-show-reservator="updateShowReservator"
                    @showDiv="$emit('showDiv', $event)"
                    @hideDiv="$emit('hideDiv')"
                />
                <ItemPeriod 
                    period-name="下午："
                    :reservator-list="schedulesList.slice(6)"
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
                    <input type="text" id="name" name="name" class="form-control" v-model="localFormInfo.title" required>
                </div>
            </div>

            <div class="row mb-2">
                <label for="email" class="col-sm-2 col-form-label">參與人：</label>
                <div class="col-sm-9">
                    <input type="text" id="email" name="email" class="form-control" v-model="localFormInfo.namesString" required>
                </div>
            </div>

            <div class="row mb-2">
                <label class="col-sm-2 col-form-label">開始時間：</label>
                <div class="col-sm-9">
                    <select class="form-select" v-model="localFormInfo.start_time" required>
                        <option :value="item" v-for="item in time_period" :key="item">{{ item }}</option>
                    </select>
                </div>
            </div>

            <div class="row mb-2">
                <label class="col-sm-2 col-form-label">結束時間：</label>
                <div class="col-sm-9">
                    <select class="form-select" v-model="localFormInfo.end_time" required>
                        <option :value="item" v-for="item in time_period" :key="item">{{ item }}</option>
                    </select>
                </div>
            </div>

            <div class="row mb-2">
                <label for="content" class="col-sm-2 col-form-label">會議內容：</label>
                <div class="col-sm-9">
                    <textarea id="content" name="content" class="form-control" maxlength="200" v-model="localFormInfo.description"></textarea>
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
  import CommWithGql from '@/components/CommWithGql.vue'

  export default {
    name: 'ReserveForm',
    emits: ['close-form', 'showDiv', 'hideDiv', 'update-form'],
    props: ['users', 'formInfo', 'formDisplay', 'roomName', 'schedulesList'],
    data() {
      return {
        time_period: [],
        showReservator: '',
      };
    },
    computed: {
        localFormInfo() {
        return JSON.parse(JSON.stringify(this.formInfo));
        },
    },
    methods: {
      submitForm() {
        
        const {dayTime, roomId } = this.formInfo;
        const startTime = dayTime+'-'+ this.localFormInfo.start_time + ':00';
        const startAt = this.transferToTimestamp(startTime);

        const endTime = dayTime+'-'+ this.localFormInfo.end_time + ':00';
        const endAt = this.transferToTimestamp(endTime);
        
        const namesArray = this.localFormInfo.namesString.split(',');
        const idsArray = namesArray.map(name => {
            const user = this.users.find(user => user.name === name);
            return user ? user.id : null;
        });
        const participantsIDs = idsArray
        const remindAt = startAt + 3600000

        const newFormInfo = { title: this.localFormInfo.title, description: this.localFormInfo.description, startAt, endAt, roomId, participantsIDs, remindAt };

        if (this.localFormInfo.eventId !== '') {
            newFormInfo.id = this.localFormInfo.eventId;
        }
        
        console.log('newFormInfo:', newFormInfo);   
        this.$refs.commWithGql.createEvent(newFormInfo);
        this.$emit('update-form');
        this.closeForm();
      },
      updateShowReservator(value) {
        this.showReservator = value;
      },
      closeForm() {
        this.$emit('close-form');
      },
      transferToTimestamp(time) {
        const formattedTime_ = time.replace(/(\d{4})-(\d{2})-(\d{2})-(\d{1,2}):(\d{2}):(\d{2})/, '$1-$2-$3T$4:$5:$6');
        const formattedTime = formattedTime_.replace(/T(\d):/, 'T0$1:');
        const date = new Date(formattedTime);
        return date.getTime();
      },
    },
    components: {
        ItemPeriod,
        CommWithGql,
    },
    mounted() {
        this.time_period = this.$names;
    }
  }
  </script>
  