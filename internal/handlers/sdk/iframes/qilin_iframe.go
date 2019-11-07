package iframes

const QilinIframe = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <script src="./qilin-store-adapter.js"></script>
  <style>
    html, body {
      height: 100%;
      overflow: hidden;
    }
  </style>
  <title>Document</title>
</head>
<body>
  <script>
    const helper = qilinGameProxy('{{.URL}}');  
    helper.init({ url: './frame.html?token=gjkgjhgijhjjh' })
      .then(() => console.log('Adapter was started'))
      .catch(err => console.log(err));
  </script>
</body>
</html>
`