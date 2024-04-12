
// $(document).ready(function() {
//   // 绑定点击事件到搜索按钮上
//   $('.main-button').click(function(e) {
//       // 阻止表单的默认提交事件
//       e.preventDefault();
      
//       // 获取表单内的值
//       var area = $('select[name="area"]').val();
//       var address = $('input[name="address"]').val();
//       var price = $('select[name="price"]').val();
      
//       // 显示弹窗
//       alert('Area: ' + area + '\nAddress: ' + address + '\nPrice: ' + price);
//   });
// });

// $(document).ready(function() {
//   $('.main-button').click(function(e) {
//       e.preventDefault();
      
//       var formData = {
//           'area': $('select[name="area"]').val(),
//           'address': $('input[name="address"]').val(),
//           'price': $('select[name="price"]').val()
//       };

//       $.ajax({
//           type: "POST",
//           url: "/submit-form",
//           data: formData,
//           success: function(response) {
//               console.log("Data Submitted: ", response);
//               alert("Data Submitted: " + response);
//           },
//           error: function(error) {
//               console.log("Error: ", error);
//           }
//       });
//   });
// });
