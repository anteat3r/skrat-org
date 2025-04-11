async function getPerm() {
  try {
  let device = await navigator.bluetooth.requestDevice({filters: [{namePrefix: "BBC"}], optionalServices: ["6e400001-b5a3-f393-e0a9-e50e24dcca9e"]});
    await device.gatt.connect();
    let service = await device.gatt.getPrimaryService("6e400001-b5a3-f393-e0a9-e50e24dcca9e");
    let characteristic = await service.getCharacteristic("6e400002-b5a3-f393-e0a9-e50e24dcca9e");
    let idk = await characteristic.startNotifications();
    document.querySelector("h1").innerText = JSON.stringify(Object.keys(idk));
  } catch (e) {
    alert(e);
  }
}
// async function getPerm() {
//   let device = await navigator.usb.requestDevice({filters: [{vendorId: 0x0d28}]});
//   console.log(device);
//   await device.open();
//   await device.selectConfiguration(1);
//   await device.claimInterface(2);
//   device.controlTransferOut({
//       requestType: 'class',
//       recipient: 'interface',
//       request: 0x22,
//       value: 0x01,
//       index: 0x02});
//   let decoder = new TextDecoder();
//   while (true) {
//     let res = await device.tranferIn(5, 5);
//     console.log(decoder.decode(res.data));
//   }
// }
