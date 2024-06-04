  <template>
    <comm-with-gql ref="commWithGql"></comm-with-gql>
  
    <div id="form-popup" v-if="formDisplay" class="container mt-3">
        <form id="submitForm">
            <div class="row mb-4">
                <div class="col-sm-12">
                    <h3 class="text-center">會議室名稱：{{ localFormInfo.roomName }}</h3>
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

            <div>
                <div class="row mb-2 align-items-center">
                    <label for="file" class="col-sm-2 col-form-label">上傳文件：</label>
                    <div class="col-sm-10 d-flex align-items-center">
                        <input type="file" id="file" name="file" class="form-control me-2 w-50" @change="onFileChange">
                        <a @click="uploadFile" class="btn btn-primary me-2" style="cursor: pointer;">上傳檔案</a>
                    </div>
                </div>
            </div>

            <div class="row mb-2">
                <label class="col-sm-2 col-form-label">Uploaded file URL:</label>
                <div class="col-sm-10">
                    <a v-if="localFormInfo.fileUrl" :href="localFormInfo.fileUrl" target="_blank">{{ localFormInfo.fileUrl }}</a>
                    <a v-else-if="fileUrl" :href="fileUrl" target="_blank">{{ fileUrl }}</a>
                </div>
            </div>

            <div class="row mb-2">
                <label class="col-sm-2 col-form-label"><strong>會議會在前十分鐘提醒</strong></label>
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
  import axios from 'axios';

  export default {
    name: 'ReserveForm',
    emits: ['close-form', 'showDiv', 'hideDiv', 'update-form'],
    props: ['users', 'formInfo', 'formDisplay', 'roomName', 'schedulesList'],
    data() {
      return {

        selectedFile: null,
        fileName: '',
        fileUrl: '',
        downLoadUrl: '',
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
      onFileChange(event) {
        this.selectedFile = event.target.files[0];
        this.localFormInfo.fileName = this.selectedFile ? this.selectedFile.name : '';
      },
      uploadFile() {
            const formData = new FormData();
            const operations = JSON.stringify({
                query: `mutation ($file: Upload!) { uploadFile(file: $file) }`,
                variables: { file: null }
            });
            const map = JSON.stringify({
                "0": ["variables.file"]
            });

            formData.append('operations', operations);
            formData.append('map', map);
            formData.append('0', this.selectedFile);

            axios.post('http://localhost:8080/query', formData, {
                headers: {
                'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImxlZWl2YW4xMDA3QGdtYWlsLmNvbSIsImV4cCI6MTcxODA2ODI5OCwibmFtZSI6Ikl2YW4gTGVlIiwic3ViIjoiNjY0NWVjZTEzNmUyYTBmMDM1OTYxYmRkIn0.u0a949cBKw2qy3uVOXikTTDGHiU5UN5eUROnpA5QHTw',
                'Content-Type': 'multipart/form-data'
                }
            }).then(response => {
                this.localFormInfo.fileUrl = response.data.data.uploadFile;
                this.fileUrl = response.data.data.uploadFile;
            }).catch(error => {
                console.error('Error uploading file:', error);
            });
      },
      submitForm() {
        
        const {dayTime, roomId} = this.formInfo;
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
        const remindAt = startAt + 600000

        const attachedFile = {
            url: this.localFormInfo.fileUrl,
            name: this.localFormInfo.fileName
        }

        const newFormInfo = {title: this.localFormInfo.title, description: this.localFormInfo.description, startAt, endAt, roomId, participantsIDs, remindAt };

        if (this.localFormInfo.eventId !== '') {
            newFormInfo.id = this.localFormInfo.eventId;
        }
        
        if (this.localFormInfo.fileUrl !== '' && this.localFormInfo.fileName !== '') {
            newFormInfo.attachedFile = attachedFile;
        }
          
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
  