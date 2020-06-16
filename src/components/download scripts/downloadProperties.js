// import pdfMake from "pdfmake/build/pdfmake";
// import pdfFonts from "pdfmake/build/vfs_fonts";
// pdfMake.vfs = pdfFonts.pdfMake.vfs;
// function download(data, name) {
//   var document = {
//     content: [
//       { text: `List of Properties in ${name}`, style: "header" },
//       {
//         style: "table",
//         table: {
//           headerRows: 1,
//           widths: ["*", 50, 60, "auto", 60, 60, "auto", 50],
//           heights: 20,
//           body: loopData(data)
//         }
//       }
//     ],
//     styles: {
//       header: {
//         fontSize: 15,
//         bold: !0,
//         margin: [0, 0, 0, 25],
//         alignment: "center",
//         decoration: "underline"
//       },
//       table: { margin: [0, 10, 10, 0] },
//       tableData: { fontSize: 11 },
//       tableHeader: { bold: !0, fontSize: 12, color: "black" }
//     },
//     defaultStyle: { alignment: "left", color: "#232323" }
//   };
//   const currentDate = new Date().toLocaleDateString("en-EN", {
//     year: "numeric",
//     month: "long",
//     day: "numeric"
//   });
//   pdfMake
//     .createPdf(document)
//     .download(`List of Properties in ${name} on ${currentDate}.pdf`);
// }
// function loopData(items) {
//   var array = [];
//   array.push([
//     { text: "Full Name", style: "tableHeader" },
//     { text: "House Code", style: "tableHeader" },
//     { text: "Phone Number", style: "tableHeader" },
//     { text: "Sector", style: "tableHeader" },
//     { text: "Cell", style: "tableHeader" },
//     { text: "Village", style: "tableHeader" },
//     { text: "Rented", style: "tableHeader" },
//     { text: "Amount", style: "tableHeader" }
//   ]);
//   items.map(item => {
//     array.push([
//       { text: `${item.owner.fname} ${item.owner.lname}`, style: "tableData" },
//       { text: item.id, style: "tableData", noWrap: !0 },
//       { text: item.owner.phone, style: "tableData", noWrap: !0 },
//       { text: item.address.sector, style: "tableData", noWrap: !0 },
//       { text: item.address.cell, style: "tableData", noWrap: !0 },
//       { text: item.address.village, style: "tableData", noWrap: !0 },
//       { text: item.occupied ? "Yes" : "No", style: "tableData" },
//       {
//         text: Number(item.due).toLocaleString() + " Rwf",
//         style: "tableData",
//         noWrap: !0
//       }
//     ]);
//   });
//   return array;
// }
// export default download;
import "jspdf-autotable";
import jsPDF from "jspdf";
const CURRENT_DATE = new Date().toLocaleDateString("en-EN", {
  year: "numeric",
  month: "long"
});
const TITLE = "Report";
const DOC_NAME = `Report of ${CURRENT_DATE}`;
const FOOTER = "Powered by PAYPACK";

function Download(param) {
  var doc = new jsPDF("l");
  var totalPagesExp = "{total_pages_count_string}";

  doc.autoTable({
    theme: "grid",
    compress: true,
    head: [param.data.Headers],
    body: param.data.Body,
    styles: { overflow: "ellipsize", cellWidth: "auto" },
    didDrawPage: function(data) {
      // Header
      doc.setFontSize(16);
      doc.setTextColor(40);
      doc.setFontStyle("normal");
      doc.text(
        param.title ? param.title : TITLE,
        data.settings.margin.left,
        22
      );
      // Footer
      var str = "Page " + doc.internal.getNumberOfPages();
      if (typeof doc.putTotalPages === "function") {
        str = str + " of " + totalPagesExp;
      }
      doc.setFontSize(10);
      var pageSize = doc.internal.pageSize;
      var pageHeight = pageSize.height ? pageSize.height : pageSize.getHeight();
      var pageWidth = pageSize.width ? pageSize.width : pageSize.getWidth();
      doc.text(str, data.settings.margin.left, pageHeight - 10);
      doc.setFontSize(11);
      doc.text(
        `${param.footer ? param.footer : FOOTER}`,
        pageWidth - data.settings.margin.left,
        pageHeight - 10,
        {
          align: "right"
        }
      );
    },
    margin: { top: 30 }
  });
  if (typeof doc.putTotalPages === "function") {
    doc.putTotalPages(totalPagesExp);
  }
  doc.save(`${param.name ? param.name : DOC_NAME}.pdf`);
}
export default Download;
