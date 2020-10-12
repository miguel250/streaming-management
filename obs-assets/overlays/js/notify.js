(async () => {
  const events = new EventSource("/events");
  const stack = []
  const audioFollowElem = document.body.getElementsByClassName("notification")[0];
  const audioSubscriberElem = document.body.getElementsByClassName("subscriber-audio")[0];
  events.addEventListener("new_follower", async (e) => {
    const obj = {
      displayName: e.data,
      eventType: "new_follower",
    }
    stack.push(obj);
  });

  events.addEventListener("new_subscriber", async (e) => {
    const obj = {
      displayName: e.data,
      eventType: "new_subscriber",
    }
    stack.push(obj);
  });


  const showNotification = () => {
    setTimeout(() => {
      const obj = stack.pop();
      let elem = null;
      let displayNameElem;
      let audioElem;

      if (obj && obj.displayName) {
        if (obj.eventType === "new_follower") {
          elem = document.body.getElementsByClassName("new-follower")[0];
          displayNameElem = document.body.getElementsByClassName("display-name")[0];
          audioElem = audioFollowElem;
        }

        if (obj.eventType === "new_subscriber") {
          elem = document.body.getElementsByClassName("new-subscriber")[0];
          displayNameElem = document.body.getElementsByClassName("sub-display-name")[0];
          audioElem = audioSubscriberElem;
        }
      }

      if (elem != null) {
        audioElem.pause();
        audioElem.volume = 1;
        audioElem.currentTime = 0;
        elem.classList.remove("show");
        displayNameElem.innerText = obj.displayName;
        elem.classList.add("show");

        const newElem = elem.cloneNode(true);
        elem.parentNode.replaceChild(newElem, elem);
        audioElem.play()
      }
      showNotification();
    }, 10000);
  };

  const loadAudio = () => {
    audioSubscriberElem.volume = 0;
    audioSubscriberElem.play().then(showNotification).catch(showNotification);
  };

  audioFollowElem.volume = 0;
  audioFollowElem.play().then(loadAudio).catch(loadAudio);
})()
