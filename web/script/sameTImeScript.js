document.addEventListener("DOMContentLoaded", function () {
    var videos = document.querySelectorAll(".video");
    var playButton = document.getElementById("playButton");
    var stopButton = document.getElementById("stopButton");

    playButton.addEventListener("click", function () {
        videos.forEach(function (video) {
            video.play();
        });
        playButton.disabled = true;
        stopButton.disabled = false;
    });

    stopButton.addEventListener("click", function () {
        videos.forEach(function (video) {
            video.pause();
            video.currentTime = 0;
        });
        playButton.disabled = false;
        stopButton.disabled = true;
    });

    // Pause videos when they end
    videos.forEach(function (video) {
        video.addEventListener("ended", function () {
            playButton.disabled = false;
            stopButton.disabled = true;
        });
    });
});
