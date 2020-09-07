const Animation = (elem, path) => {
  return bodymovin.loadAnimation({
    container: elem,
    renderer: 'svg',
    loop: true,
    autoplay: true,
    path: path,
  });
};
