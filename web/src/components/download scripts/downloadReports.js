import "jspdf-autotable";
import jsPDF from "jspdf";
const CURRENT_DATE = new Date().toLocaleDateString("en-EN", {
  year: "numeric",
  month: "long"
});
const TITLE = "Monthly report";
const DOC_NAME = `Monthly Report of ${CURRENT_DATE}`;
const FOOTER = "Powered by PAYPACK";
const SHOWHEAD = "everyPage";

function Download(param) {
  var doc = new jsPDF({ orientation: "portrait", format: "a4" });
  var totalPagesExp = "{total_pages_count_string}";
  param.data.forEach(item => {
    doc.autoTable({
      theme: "grid",
      compress: true,
      showHead: item.SHOWHEAD ? item.SHOWHEAD : SHOWHEAD,
      headStyles: { fillColor: [1, 125, 179] },
      body: item.BODY,
      columns: item.COLUMNS,
      styles: { overflow: "ellipsize", cellWidth: "auto" },
      margin: { top: 30 },
      didDrawPage: function(data) {
        // Header
        doc.setFontSize(16);
        doc.setTextColor(40);
        doc.setFontStyle("normal");
        doc.text(
          param.config.TITLE ? param.config.TITLE : TITLE,
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
        var pageHeight = pageSize.height
          ? pageSize.height
          : pageSize.getHeight();
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
      }
    });
  });
  if (typeof doc.putTotalPages === "function") {
    doc.putTotalPages(totalPagesExp);
  }
  doc.save(`${param.config.name ? param.config.name : DOC_NAME}.pdf`);
}
export default Download;
