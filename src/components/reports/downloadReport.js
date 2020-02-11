import pdfMake from "pdfmake/build/pdfmake";
import pdfFonts from "pdfmake/build/vfs_fonts";
pdfMake.vfs = pdfFonts.pdfMake.vfs;

function download(user, payment) {
  const name = user.owner.fname + " " + user.owner.lname;
  var document = {
    content: [
      {
        text: `Details of ${name}:`,
        style: "header"
      },
      {
        style: "table",
        table: {
          headerRows: 0,
          widths: ["*", "*"],
          body: [
            [
              {
                text: "Names",
                style: "tableHeader"
              },
              {
                text: name,
                style: "tableData"
              }
            ],
            [
              {
                text: "Phone Number",
                style: "tableHeader"
              },
              {
                text: user.owner.phone,
                style: "tableData"
              }
            ],
            [
              {
                text: "House ID",
                style: "tableHeader"
              },
              {
                text: user.id,
                style: "tableData"
              }
            ],
            [
              {
                text: "Location",
                style: "tableHeader"
              },
              {
                text:
                  user.address.sector +
                  ", " +
                  user.address.cell +
                  ", " +
                  user.address.village,
                style: "tableData"
              }
            ],
            [
              {
                text: "Amount",
                style: "tableHeader"
              },
              {
                text: Number(user.due).toLocaleString() + " Rwf",
                style: "tableData"
              }
            ],
            [
              {
                text: "For Rent",
                style: "tableHeader"
              },
              {
                text: user.occupied ? "Yes" : "No",
                style: "tableData"
              }
            ],
            [
              {
                text: "Registered by",
                style: "tableHeader"
              },
              {
                text: user.recorded_by,
                style: "tableData"
              }
            ],
            [
              {
                text: "Registered on",
                style: "tableHeader"
              },
              {
                text: new Date(user.created_at).toLocaleString("en-EN", {
                  year: "numeric",
                  month: "long",
                  day: "numeric"
                }),
                style: "tableData"
              }
            ]
          ]
        }
      },
      {
        text: `Payment History of ${name}:`,
        style: "header"
      },
      {
        style: "table",
        table: {
          headerRows: 1,
          widths: ["*", "*"],
          body: BodyData(payment)
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
  pdfMake.createPdf(document).download(`${name} Report.pdf`);
}

function BodyData(items) {
  let data = [];
  data.push([
    {
      text: "Months",
      style: "tableHeader"
    },
    {
      text: "Status",
      style: "tableHeader"
    }
  ]);
  items.map(item => {
    data.push([
      {
        text: new Date(item.created_at).toLocaleString("en-EN", {
          month: "long"
        }),
        style: "tableData"
      },
      {
        text: item.status == "pending" ? "Not Payed" : "Payed",
        style: "tableData"
      }
    ]);
  });
  return data;
}

export default download;
