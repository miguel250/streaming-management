(() => {
  const events = new EventSource("/events");
  const stack = []

  events.addEventListener("new_chat_message", async (e) => {
    stack.push(JSON.parse(e.data))
  })

  const showChatMessage = () => {
    setTimeout(() => {
      const data = stack.pop();

      if (!data) {
        showChatMessage();
        return
      }

      const box = document.createElement("div");
      box.classList.add('box')

      if (data.profile_image !== "") {
        const profileImage = document.createElement("div");
        const profileImg = document.createElement("img");

        profileImg.src = data.profile_image;
        profileImage.appendChild(profileImg);
        profileImage.classList.add('profile-image');
        box.appendChild(profileImage);
      }

      const message = document.createElement("div");
      message.classList.add('message')

      message.innerHTML = `${data.message}`;

      const username = document.createElement('div');

      const usernameSpan = document.createElement('span');
      usernameSpan.innerText = `${data['display-name']}`;
      username.appendChild(usernameSpan);

      if (data.badges != null) {
        data.badges.forEach((badge) => {
          const img = document.createElement('img');
          img.src = badge.image_url_2x;
          img.classList.add("badge");
          username.appendChild(img);
        });
      }

      username.classList.add('username');
      box.appendChild(message)
      box.appendChild(username)
      document.body.appendChild(box);

      setTimeout(() => {
        box.classList.add('box-hide')
        box.addEventListener('transitionend', () => {
          try {
            document.body.removeChild(box);
          } catch (e) { }
        });
      }, 10000);
      showChatMessage();
    }, 800);
  };

  const pageScroll = () => {
    window.scrollBy(0, 1);
    scrolldelay = setTimeout(pageScroll, 10);
  }
  pageScroll();
  showChatMessage();
})()
