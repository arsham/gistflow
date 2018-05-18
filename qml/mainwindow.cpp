#include "mainwindow.h"
#include "ui_mainwindow.h"
#include <QUiLoader>
#include <QFile>

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    QUiLoader loader;

    QFile file(":/ui/settings.ui");
    file.open(QFile::ReadOnly);
    QWidget *settings = loader.load(&file, this);
    file.close();

    file.setFileName(":/ui/tabheader.ui");
    file.open(QFile::ReadOnly);
    QWidget *header = loader.load(&file, this);
    file.close();

    file.setFileName(":/ui/gist.ui");
    file.open(QFile::ReadOnly);
    QWidget *g= loader.load(&file, this);
    file.close();

    ui->tabWidget->addTab(header, "Gist");
    QBoxLayout* layout = header->findChild<QVBoxLayout *>("mainLayout", Qt::FindChildrenRecursively);
    layout->addWidget(g);
    ui->tabWidget->addTab(settings, "Settings");
}

MainWindow::~MainWindow()
{
    delete ui;
}
