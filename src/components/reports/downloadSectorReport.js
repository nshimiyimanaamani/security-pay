import pdfMake from "pdfmake/build/pdfmake";
import pdfFonts from "pdfmake/build/vfs_fonts";
pdfMake.vfs = pdfFonts.pdfMake.vfs;

function download(sector, cell, sectorName, date) {
  var document = {
    content: [
      {
        text: `Report of ${sectorName}:`,
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
                text: sector.total,
                style: "tableHeader"
              },
              {
                text: sector.payed,
                style: "tableData"
              },
              {
                text: Number(sector.payedAmount).toLocaleString() + " Rwf",
                style: "tableData"
              },
              {
                text: sector.pending,
                style: "tableData"
              },
              {
                text: Number(sector.unpayedAmount).toLocaleString() + " Rwf",
                style: "tableData"
              }
            ]
          ]
        }
      },
      {
        text: `${sectorName} cells:`,
        style: "header"
      },
      {
        style: "table",
        table: {
          headerRows: 1,
          widths: ["*", "*", "*", "*", "*", "*"],
          body: getData(cell)
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
  pdfMake.createPdf(document).download(`${sectorName} Report of ${date}.pdf`);
}

function getData(items) {
  var array = [];
  array.push([
    {
      text: "Cell",
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
