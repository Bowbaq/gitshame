hljs.initHighlightingOnLoad();

$(document).ready(function() {
  var form = $("#submit-shame");
  var url = form.find(".link-input");

  form.submit(function(e){
    e.preventDefault();

    $.ajax({
      type: "POST",
      url: "/shame",
      contentType: "application/json",
      data: JSON.stringify({
        "URL": url.val(),
      }),
      success: function(result){
        console.log(result);
      }
    });
  });
});
