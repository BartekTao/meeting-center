  <template>
    <comm-with-gql ref="commWithGql"></comm-with-gql>
    <div id="form-popup" v-if="commentDisplay" class="container mt-3">
        <form id="submitForm">
            <div class="row mb-4">
                <div class="col-sm-12">
                    <h3 class="text-center">編輯結論</h3>
                </div>
            </div>
            <div class="row mb-2">
                <label for="name" class="col-sm-2 col-form-label">會議標題：</label>
                <div class="col-sm-9">
                    <input type="text" id="name" name="name" class="form-control" v-model="localFormInfo.title" disabled>
                </div>
            </div>

            <div class="row mb-2 comment-text">
                <label for="content" class="col-sm-2">會議結論：</label>
                <div class="col-sm-9">
                    <textarea id="content" name="content" class="form-control" maxlength="500" v-model="localFormInfo.summary"  style="height: 200px;"></textarea>
                </div>
            </div>

            <div class="row mb-2">
                <div class="col-sm-12 d-flex justify-content-center">
                    <button type="submit" class="btn btn-primary margin-right-2cm" @click.prevent="submitForm">編輯</button>
                    <button type="button" class="btn btn-secondary" @click="closeCommentForm">取消</button>
                </div>
            </div>
        </form>
    </div>
  </template>
  
  <script>
  import CommWithGql from '@/components/CommWithGql.vue'

  export default {
    name: 'ReserveForm',
    emits: ['close-comment-form', 'update-form'],
    props: ['users','formInfo', 'roomInfo', 'commentDisplay', 'userName'],
    components: {
        CommWithGql,
    },
    data() {
      return {
        time_period: [],
      };
    },
    computed: {
        localFormInfo() {
        return JSON.parse(JSON.stringify(this.formInfo));
        },
    },
    methods: {
      submitForm() {
        
        const newFormInfo = { summary: this.localFormInfo.summary, id: this.localFormInfo.eventId};

        this.$refs.commWithGql.editSummary(newFormInfo);
        this.$emit('update-form');
        this.closeCommentForm();
      },
      closeCommentForm() {
        this.$emit('close-comment-form');
      },
      transferToTimestamp(time) {
        const formattedTime_ = time.replace(/(\d{4})-(\d{2})-(\d{2})-(\d{1,2}):(\d{2}):(\d{2})/, '$1-$2-$3T$4:$5:$6');
        const formattedTime = formattedTime_.replace(/T(\d):/, 'T0$1:');
        const date = new Date(formattedTime);
        return date.getTime();
      },
    },
    mounted() {
      this.time_period = this.$names;
    }
  }
  </script>
