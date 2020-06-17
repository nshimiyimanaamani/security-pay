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
