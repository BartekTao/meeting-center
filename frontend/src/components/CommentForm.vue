  <template>
    <div id="form-popup" v-if="commentDisplay" class="container mt-3">
        <form id="submitForm">
            <div class="row mb-4">
                <div class="col-sm-12">
                    <h3 class="text-center">會議室名稱：{{ roomInfo.roomName }}</h3>
                </div>
            </div>
            <div class="row mb-2">
                <label for="name" class="col-sm-2 col-form-label">會議標題：</label>
                <div class="col-sm-9">
                    <input type="text" id="name" name="name" class="form-control" v-model="formInfo.name" disabled>
                </div>
            </div>

            <div class="row mb-2 comment-text">
                <label for="content" class="col-sm-2">會議結論：</label>
                <div class="col-sm-9">
                    <textarea id="content" name="content" class="form-control" maxlength="500" v-model="formInfo.content"  style="height: 200px;"></textarea>
                </div>
            </div>

            <div class="row mb-2">
                <div class="col-sm-12 d-flex justify-content-center">
                    <button type="submit" class="btn btn-primary margin-right-2cm" @click.prevent="submitForm">預約</button>
                    <button type="button" class="btn btn-secondary" @click="closeCommentForm">取消</button>
                </div>
            </div>
        </form>
    </div>
  </template>
  
  <script>

  export default {
    name: 'ReserveForm',
    emits: ['close-comment-form'],
    props: ['roomInfo', 'commentDisplay', 'userName'],
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
        
        this.closeCommentForm();
      },
      updateShowReservator(value) {
        this.showReservator = value;
      },
      closeCommentForm() {
        this.$emit('close-comment-form');
      }
    },
    mounted() {
      this.time_period = this.$names;
    }
  }
  </script>
