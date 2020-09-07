(() => {
  const urlQueryParams = new URLSearchParams(window.location.search);
  const startingMinutes = urlQueryParams.get('min_starting') || 10;
  let secondsLeft = +startingMinutes * 60;

  const timeInterval = setInterval(() => {
    secondsLeft--;

    const { minutes, seconds } = format(secondsLeft);
    updateElements(minutes, seconds);
    if (secondsLeft == 0) {
      clearInterval(timeInterval);
      removePop(minutesElem);
      removePop(secondsElem);
    }
  }, 1000);

  const format = (timeLeft) => {
    let minutes = Math.floor(timeLeft / 60);
    let seconds = timeLeft % 60;

    if (seconds < 10) {
      seconds = `0${seconds}`;
    }

    if (minutes < 10) {
      minutes = `0${minutes}`
    }
    return { minutes, seconds }
  };

  const parentElem = document.body.getElementsByClassName('timer')[0];
  const minutesElem = parentElem.getElementsByClassName('minutes')[0];
  const secondsElem = parentElem.getElementsByClassName('seconds')[0];

  const updateElements = (minutes, seconds) => {
    if (minutesElem.innerHTML !== minutes) {
      minutesElem.textContent = minutes;
      addPop(minutesElem);
    } else {
      removePop(minutesElem);
    }

    if (secondsElem.innerHTML !== seconds) {
      secondsElem.textContent = seconds;
      addPop(secondsElem);
    } else {
      removePop(secondsElem);
    }
  };


  const addPop = (elem) => {
    elem.classList.add('pop');
  };

  const removePop = (elem) => {
    elem.classList.remove('pop');
  };

  updateElements(startingMinutes, '00');

})();
