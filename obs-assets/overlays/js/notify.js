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
      const obj = stack.shift();
      let elem = null;
      let displayNameElem;

      if (!obj) {
        showNotification();
        return
      }

      if (obj.eventType === "new_follower") {
        elem = document.body.getElementsByClassName("new-follower")[0];
        displayNameElem = document.body.getElementsByClassName("display-name")[0];
      }

      if (obj.eventType === "new_subscriber") {
        elem = document.body.getElementsByClassName("new-subscriber")[0];
        displayNameElem = document.body.getElementsByClassName("sub-display-name")[0];
      }

      if (elem != null) {
        elem.classList.remove("show");
        displayNameElem.innerText = obj.displayName;
        elem.classList.add("show");

        const newElem = elem.cloneNode(true);
        elem.parentNode.replaceChild(newElem, elem);
        if (obj.eventType === "new_subscriber") {
          audioSubscriberElem.currentTime = 0;
          audioSubscriberElem.volume = 1;
          audioSubscriberElem.play().then().catch(() => {
            audioSubscriberElem.play();
          });
        }

        if (obj.eventType === "new_follower") {
          audioFollowElem.currentTime = 0;
          audioFollowElem.volume = 1;
          audioFollowElem.play().then().catch(() => {
            audioFollowElem.play();
          });
        }
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
