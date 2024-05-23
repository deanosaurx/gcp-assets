document.querySelectorAll("a[data-asset-type]").forEach((link) => {
  link.addEventListener("click", function (event) {
    event.preventDefault();
    const assetType = this.getAttribute("data-asset-type");
    fetch(`/assets?type=compute.googleapis.com/${assetType}`)
      .then((response) => response.json())
      .then((data) => {
        const contentDiv = document.getElementById("content");
        contentDiv.innerHTML = renderTable(data);
      })
      .catch((error) => console.error("Error fetching assets:", error));
  });
});

function renderTable(assets) {
  let table =
    "<table><thead><tr><th>Name</th><th>Asset Type</th><th>Project</th><th>Resource ID</th></tr></thead><tbody>";
  assets.forEach((asset) => {
    table += `<tr><td>${asset.name}</td><td>${asset.asset_type}</td><td>${asset.project}</td><td>${asset.resource_id}</td></tr>`;
  });
  table += "</tbody></table>";
  return table;
}
