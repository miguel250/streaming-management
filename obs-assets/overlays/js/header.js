(async () => {
  const getFollowerData = async () => {
    const response = await fetch("/api/goals");
    return response.json();
  };

  const updateHeader = async () => {
    const data = await getFollowerData();
    const followerNameElem = document.body.getElementsByClassName("new_follower_name")[0];
    const followerCounterElem = document.body.getElementsByClassName("follower_counter")[0];
    followerNameElem.innerText = data.follower_name;
    if (data.disable_follower_goal) {
      document.body.getElementsByClassName("follower-goal")[0].remove();
      return
    }
    followerCounterElem.innerText = `${data.follower_total} / ${data.follower_goal}`

    const subscriberNameElem = document.body.getElementsByClassName("new_subscriber_name")[0];
    const subscriberCounterElem = document.body.getElementsByClassName("subscriber_counter")[0];

    subscriberNameElem.innerText = data.subscriber_name
    subscriberCounterElem.innerText = `${data.subscriber_total} / ${data.subscriber_goal}`
  }

  const events = new EventSource("/events");
  events.addEventListener("new_follower", async (e) => {
    await updateHeader();
  });

  events.addEventListener("new_subscriber", async (e) => {
    await updateHeader();
  });

  await updateHeader();
})();
