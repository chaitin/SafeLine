

$(document).ready(function () {
    var myCarousel = document.getElementById('carouselExampleCaptions')
    console.log(myCarousel)
    if (myCarousel) {
        var carousel = new bootstrap.Carousel(myCarousel, {
            interval: 2000,
            wrap: false
        })
    }
});