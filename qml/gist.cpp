#include "gist.h"
#include "ui_gist.h"

Gist::Gist(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::Gist)
{
    ui->setupUi(this);
}

Gist::~Gist()
{
    delete ui;
}
