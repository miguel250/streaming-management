(async () => {
  const getFollowerData = async () => {
    const response = await fetch("/api/goals");
    return response.json();
  };

  const updateHeader = async () => {
    const data = await getFollowerData();
    const followerNameElem = document.body.getElementsByClassName("new_follower_name")[0];
    const counterElem = document.body.getElementsByClassName("counter")[0];
    followerNameElem.innerText = data.follower_name;
    if (data.disable_follower_goal) {
      document.body.getElementsByClassName("follower-goal")[0].remove();
      return
    }
    counterElem.innerText = `${data.total} / ${data.goal}`
  }

  const events = new EventSource("/events");
  events.addEventListener("new_follower", async (e) => {
    await updateHeader();
  });
  await updateHeader();
})();
