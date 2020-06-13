import pdfMake from "pdfmake/build/pdfmake";
import pdfFonts from "pdfmake/build/vfs_fonts";
pdfMake.vfs = pdfFonts.pdfMake.vfs;

function download(cell, village, cellName, date) {
  var document = {
    content: [
      {
        text: `Report of ${cellName}:`,
        style: "header"
      },
      {
        style: "table",
        table: {
          headerRows: 1,
          widths: ["*", "*", "*", "*", "*"],
          body: [
            [
              {
                text: "No of Houses",
                style: "tableHeader"
              },
              {
                text: "Paid Houses",
                style: "tableHeader"
              },
              {
                text: "Paid Amount",
                style: "tableHeader"
              },
              {
                text: "unPaid Houses",
                style: "tableHeader"
              },
              {
                text: "unPaid Amount",
                style: "tableHeader"
              }
            ],
            [
              {
                text: cell.total,
                style: "tableHeader"
              },
              {
                text: cell.payed,
                style: "tableData"
              },
              {
                text: Number(cell.payedAmount).toLocaleString() + " Rwf",
                style: "tableData"
              },
              {
                text: cell.pending,
                style: "tableData"
              },
              {
                text: Number(cell.unpayedAmount).toLocaleString() + " Rwf",
                style: "tableData"
              }
            ]
          ]
        }
      },
      {
        text: `${cellName} villages:`,
        style: "header"
      },
      {
        style: "table",
        table: {
          headerRows: 1,
          widths: ["*", "*", "*", "*", "*", "*"],
          body: getData(village)
        }
      }
    ],
    styles: {
      header: {
        fontSize: 15,
        bold: true,
        margin: [0, 0, 0, 25],
        alignment: "center",
        decoration: "underline"
      },
      table: {
        margin: [0, 10, 10, 0]
      },
      tableData: {
        fontSize: 11
      },
      tableHeader: {
        bold: true,
        fontSize: 12,
        color: "black"
      }
    },
    defaultStyle: {
      alignment: "left",
      color: "#232323"
    }
  };
  pdfMake.createPdf(document).download(`${cellName} Report of ${date}.pdf`);
}

function getData(items) {
  var array = [];
  array.push([
    {
      text: "Village",
      style: "tableHeader"
    },
    {
      text: "No of Houses",
      style: "tableHeader"
    },
    {
      text: "Paid Houses",
      style: "tableHeader"
    },
    {
      text: "Paid Amount",
      style: "tableHeader"
    },
    {
      text: "unPaid Houses",
      style: "tableHeader"
    },
    {
      text: "unPaid Amount",
      style: "tableHeader"
    }
  ]);
  items.map(item => {
    array.push([
      {
        text: item.name,
        style: "tableData"
      },
      {
        text: item.total,
        style: "tableData"
      },
      {
        text: item.payed,
        style: "tableData"
      },
      {
        text: Number(item.payedAmount).toLocaleString() + " Rwf",
        style: "tableData"
      },
      {
        text: item.pending,
        style: "tableData"
      },
      {
        text: Number(item.unpayedAmount).toLocaleString() + " Rwf",
        style: "tableData"
      }
    ]);
  });
  return array;
}
export default download;
