export default name => {
  if (!name) {
    document.title = "Paypack";
    return;
  }
  document.title = `Paypack | ${name}`;
};
