import pdfMake from "pdfmake/build/pdfmake";
import pdfFonts from "pdfmake/build/vfs_fonts";
pdfMake.vfs = pdfFonts.pdfMake.vfs;

function download(data, village) {
  var document = {
    content: [
      {
        text: `Report of ${village}:`,
        style: "header"
      },
      {
        style: "table",
        table: {
          headerRows: 1,
          widths: ["*", "*", "*", "*", "*"],
          body: [
            [
              { text: "No of Houses", style: "tableHeader" },
              { text: "No of Payed Houses", style: "tableHeader" },
              { text: "Payed Amount", style: "tableHeader" },
              { text: "No of unpayed Houses", style: "tableHeader" },
              { text: "UnPayed Amount", style: "tableHeader" }
            ],
            [
              { text: data.total, style: "tableHeader" },
              { text: data.payed, style: "tableData" },
              {
                text: Number(data.payedAmount).toLocaleString() + " Rwf",
                style: "tableData"
              },
              { text: data.pending, style: "tableData" },
              {
                text: Number(data.unpayedAmount).toLocaleString() + " Rwf",
                style: "tableData"
              }
            ]
          ]
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
  pdfMake.createPdf(document).download(`${village} Report.pdf`);
}

export default download;
