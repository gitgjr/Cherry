document.addEventListener("DOMContentLoaded", function () {
    var dataFolderPath = "../data/"
    var videoList = document.getElementById("videoList");
    var videoPlayer = document.getElementById("videoPlayer");

    // Replace these with your actual video file names
    var videoFiles = ["earth.mp4", "video2.mp4", "video3.mp4"];

    videoFiles.forEach(function (file) {
        var fileName = file.split("/").pop().split(".")[0];
        var listItem = document.createElement("li");
        listItem.textContent = fileName;
        let filePath = dataFolderPath+file
        listItem.setAttribute("data-src", filePath);
        videoList.appendChild(listItem);
    });

    videoList.addEventListener("click", function (event) {
        if (event.target.tagName === "LI") {
            var videoSource = event.target.getAttribute("data-src");
            videoPlayer.src = videoSource;
            videoPlayer.load();
            videoPlayer.play();
        }
    });
});



