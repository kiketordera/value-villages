$('.image_cover').each(function(){
var imageWidth = $(this).find('img').width();
var imageheight = $(this).find('img'). height();
  if(imageWidth > imageheight){
    $(this).addClass('landscape');
  }else{
    $(this).addClass('potrait');
  }
})

$(document).ready(function () {
	$('.main-image').flickity({
            //Options
            cellAlign: 'left',
            contain: true
     });
});
