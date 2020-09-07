(async () => {
  const events = new EventSource("/events");
  const stack = []
  const audioElem = document.body.getElementsByClassName("notification")[0];
  events.addEventListener("new_follower", async (e) => {
    stack.push(e.data);
  });

  const showNotification = () => {
    setTimeout(() => {
      const displayName = stack.pop();

      if (displayName) {
        const elem = document.body.getElementsByClassName("new-follower")[0];
        const displayNameElem = document.body.getElementsByClassName("display-name")[0];

        audioElem.pause();
        audioElem.volume = 1;
        audioElem.currentTime = 0;
        elem.classList.remove("show");
        displayNameElem.innerText = displayName;
        elem.classList.add("show");

        const newElem = elem.cloneNode(true);
        elem.parentNode.replaceChild(newElem, elem);
        audioElem.play()
      }
      showNotification();
    }, 10000);
  };

  audioElem.volume = 0
  audioElem.play().then(showNotification).catch(showNotification);
})()
